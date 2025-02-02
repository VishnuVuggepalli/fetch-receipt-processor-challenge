package repository

import (
	"sync"

	"github.com/google/uuid"
)

type ReceiptRepository interface {
	Store(points int) string
	Retrieve(id string) (int, bool)
}

type inMemoryReceiptRepository struct {
	mu    sync.RWMutex
	store map[string]int
}

func NewInMemoryReceiptRepository() ReceiptRepository {
	return &inMemoryReceiptRepository{
		store: make(map[string]int),
	}
}

func (r *inMemoryReceiptRepository) Store(points int) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	r.store[id] = points
	return id
}

func (r *inMemoryReceiptRepository) Retrieve(id string) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	points, exists := r.store[id]
	return points, exists
}
