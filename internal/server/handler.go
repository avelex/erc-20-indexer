package server

import (
	"context"

	"github.com/avelex/erc-20-indexer/internal/abi"
	"github.com/avelex/erc-20-indexer/internal/repository/memory"
)

type Handler struct {
	repo *memory.Repository
}

func New(repo *memory.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) LastEvents(ctx context.Context, limit int) ([]*abi.Erc20Transfer, error) {
	return h.repo.LastEvents(ctx, limit)
}
