# Receipt Processor Service

## Quick Start

### Using Docker
```bash
docker build -t receipt-refactored .
docker run -p 8080:8080 receipt-refactored 
```

### Local Development
```bash
go mod init 
go mod tidy
go run cmd/fetch-receipt-processor-challenge/main.go
```

## API Endpoints

### Process Receipt
```bash
# POST /receipts/process
curl -X POST \
  -H "Content-Type: application/json" \
  -d @examples/simple-receipt.json \
  http://localhost:8080/receipts/process
```

### Get Points
```bash
# GET /receipts/{id}/points
curl http://localhost:8080/receipts/{id}/points
```

## Project Structure
```
.
├── cmd/
│   └── fetch-receipt-processor-challenge/
│       └── main.go
├── internal/
│   ├── controller/
│   ├── model/
│   ├── repository/
│   └── service/
├── examples/
│   ├── simple-receipt.json
│   └── morning-receipt.json
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Development

### Prerequisites
- Go 1.23+
- Docker (optional)

## API Documentation

Detailed API documentation is available in `api.yml` (OpenAPI 3.0 specification).

## Points Calculation Rules

Points are awarded based on the following criteria:
1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items
5. Bonus points for item descriptions (length multiple of 3)
6. 6 points for odd purchase dates
7. 10 points for purchases between 2:00 PM and 4:00 PM