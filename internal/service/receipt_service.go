package service

import (
	"errors"
	"fetch-receipt-processor-challenge/internal/model"
	"fetch-receipt-processor-challenge/internal/repository"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	ErrInvalidReceipt  = errors.New("invalid receipt data")
	ErrInvalidTotal    = errors.New("invalid total format")
	ErrInvalidPrice    = errors.New("invalid price format")
	ErrInvalidDate     = errors.New("invalid purchase date")
	ErrInvalidTime     = errors.New("invalid purchase time")
	ErrInvalidItem     = errors.New("invalid item data")
	ErrNotFound        = errors.New("receipt not found")
)

type ReceiptService interface {
	ProcessReceipt(receipt *model.Receipt) (string, error)
	GetPoints(id string) (int, error)
}

type receiptServiceImpl struct {
	repo repository.ReceiptRepository
}

func NewReceiptService(repo repository.ReceiptRepository) ReceiptService {
	return &receiptServiceImpl{repo: repo}
}

func (s *receiptServiceImpl) ProcessReceipt(receipt *model.Receipt) (string, error) {
	if err := validateReceipt(receipt); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidReceipt, err)
	}

	points, err := calculatePoints(receipt)
	if err != nil {
		return "", fmt.Errorf("points calculation failed: %w", err)
	}

	return s.repo.Store(points), nil
}

func (s *receiptServiceImpl) GetPoints(id string) (int, error) {
	points, exists := s.repo.Retrieve(id)
	if !exists {
		return 0, ErrNotFound
	}
	return points, nil
}

func validateReceipt(receipt *model.Receipt) error {
	if strings.TrimSpace(receipt.Retailer) == "" {
		return fmt.Errorf("retailer name cannot be empty")
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDate, err)
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTime, err)
	}

	if len(receipt.Items) == 0 {
		return errors.New("receipt must contain at least one item")
	}

	if _, err := parseTotal(receipt.Total); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTotal, err)
	}

	for i, item := range receipt.Items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			return fmt.Errorf("%w: item %d has empty description", ErrInvalidItem, i+1)
		}
		if _, err := parsePrice(item.Price); err != nil {
			return fmt.Errorf("%w: item %d: %v", ErrInvalidPrice, i+1, err)
		}
	}

	return nil
}

func calculatePoints(receipt *model.Receipt) (int, error) {
	points := 0

	// 1. Retailer name alphanumeric characters
	points += countAlphanumeric(receipt.Retailer)

	// 2. Round dollar amount
	totalCents, err := parseTotal(receipt.Total)
	if err != nil {
		return 0, err
	}
	if totalCents%100 == 0 {
		points += 50
	}

	// 3. Multiple of 0.25
	if totalCents%25 == 0 {
		points += 25
	}

	// 4. Items pairs
	points += (len(receipt.Items) / 2) * 5

	// 5. Item description length multiple of 3
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			priceCents, err := parsePrice(item.Price)
			if err != nil {
				return 0, err
			}
			priceDollars := float64(priceCents) / 100.0
			points += int(math.Ceil(priceDollars * 0.2))
		}
	}

	// 6. Odd purchase day
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// 7. Purchase time between 2:00pm and 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if isBetweenTwoAndFourPM(purchaseTime) {
		points += 10
	}

	return points, nil
}

func countAlphanumeric(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

func parseTotal(totalStr string) (int, error) {
	validTotal := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !validTotal.MatchString(totalStr) {
		return 0, ErrInvalidTotal
	}

	parts := strings.Split(totalStr, ".")
	dollars, _ := strconv.Atoi(parts[0])
	cents, _ := strconv.Atoi(parts[1])
	return dollars*100 + cents, nil
}

func parsePrice(priceStr string) (int, error) {
	validPrice := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !validPrice.MatchString(priceStr) {
		return 0, ErrInvalidPrice
	}

	parts := strings.Split(priceStr, ".")
	dollars, _ := strconv.Atoi(parts[0])
	cents, _ := strconv.Atoi(parts[1])
	return dollars*100 + cents, nil
}

func isBetweenTwoAndFourPM(t time.Time) bool {
	start := time.Date(t.Year(), t.Month(), t.Day(), 14, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 16, 0, 0, 0, t.Location())
	return t.After(start) && t.Before(end)
}