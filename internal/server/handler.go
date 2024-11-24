package server

import (
	"context"

	"github.com/avelex/erc-20-indexer/internal/queries"
	"github.com/avelex/erc-20-indexer/internal/repository/timescale"
)

type Handler struct {
	repo *timescale.Repository
}

func New(repo *timescale.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) LastEvents(ctx context.Context, limit int) ([]queries.Erc20Transfer, error) {
	return h.repo.LastEvents(ctx, limit)
}
