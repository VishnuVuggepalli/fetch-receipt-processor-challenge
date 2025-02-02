// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func main() {

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("POST /receipts/process", controller.processReceipt)
// 	mux.HandleFunc("GET /receipts/{id}/points", server.getPoints)

// 	fmt.Println("Starting server on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", mux))
// }
// // curl http://localhost:8080/receipts/0bbb8ed3-12b8-494b-a7e7-42fad71a9c35/points

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fetch-receipt-processor-challenge/internal/controller"
	"fetch-receipt-processor-challenge/internal/repository"
	"fetch-receipt-processor-challenge/internal/service"
)

func main() {
	// Initialize dependencies
	repo := repository.NewInMemoryReceiptRepository()
	receiptService := service.NewReceiptService(repo)
	receiptController := controller.NewReceiptController(receiptService)

	// Configure router with middleware
	router := http.NewServeMux()
	router.HandleFunc("POST /receipts/process", receiptController.ProcessReceipt)
	router.HandleFunc("GET /receipts/{id}/points", receiptController.GetPoints)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         getPort(),
		Handler:      withLogging(withTimeout(router)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Graceful shutdown setup
	serverCtx, serverStop := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		serverStop()
	}()

	// Start server
	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Wait for server context cancellation
	<-serverCtx.Done()
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return ":" + port
}

// Middleware chain
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

func withTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
