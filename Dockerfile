FROM golang:1.23-alpine

WORKDIR /app

# Copy all source files
COPY . .

# Initialize module if needed and build
RUN if [ ! -f go.mod ]; then go mod init fetch-receipt-processor-challenge; fi && \
    go mod tidy && \
    go build -o /fetch-receipt-processor-challenge ./cmd/fetch-receipt-processor-challenge

EXPOSE 8080

CMD ["/fetch-receipt-processor-challenge"]

