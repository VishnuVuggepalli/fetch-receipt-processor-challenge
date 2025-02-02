package controller

import (
	"encoding/json"
	"errors"
	"fetch-receipt-processor-challenge/internal/model"
	"fetch-receipt-processor-challenge/internal/service"
	"net/http"
)

type ReceiptController struct {
	service service.ReceiptService
}

func NewReceiptController(s service.ReceiptService) *ReceiptController {
	return &ReceiptController{service: s}
}

// ProcessReceipt handles POST /receipts/process
func (c *ReceiptController) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	id, err := c.service.ProcessReceipt(&receipt)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, model.IdResponse{ID: id})
}

// GetPoints handles GET /receipts/{id}/points
func (c *ReceiptController) GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing receipt ID")
		return
	}

	points, err := c.service.GetPoints(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, model.PointsResponse{Points: points})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidReceipt):
		respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound):
		respondWithError(w, http.StatusNotFound, "Receipt not found")
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}
