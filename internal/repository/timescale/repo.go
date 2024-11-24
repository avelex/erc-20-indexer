package timescale

import (
	"context"
	"time"

	"github.com/avelex/erc-20-indexer/internal/abi"
	"github.com/avelex/erc-20-indexer/internal/queries"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	queries *queries.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		queries: queries.New(db),
	}
}

func (r *Repository) SaveEvent(ctx context.Context, event *abi.Erc20Transfer) error {
	return r.queries.SaveTransfer(ctx, queries.SaveTransferParams{
		Timestamp:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Sender:       event.From.String(),
		Recipient:    event.To.String(),
		Amount:       pgtype.Numeric{Int: event.Value, Valid: true},
		TxHash:       event.Raw.TxHash.String(),
		TokenAddress: event.Raw.Address.String(),
	})
}

func (r *Repository) LastEvents(ctx context.Context, limit int) ([]queries.Erc20Transfer, error) {
	return r.queries.GetLastEvents(ctx, int32(limit))
}
