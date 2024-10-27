package memory

import (
	"context"
	"sync"

	"github.com/avelex/erc-20-indexer/internal/abi"
)

type Repository struct {
	mu     *sync.RWMutex
	events []*abi.Erc20Transfer
}

func New() *Repository {
	return &Repository{
		events: make([]*abi.Erc20Transfer, 0),
		mu:     &sync.RWMutex{},
	}
}

func (r *Repository) SaveEvent(ctx context.Context, event *abi.Erc20Transfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)

	return nil
}

func (r *Repository) LastEvents(ctx context.Context, limit int) ([]*abi.Erc20Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.events[len(r.events)-limit:], nil
}
