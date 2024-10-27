package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/avelex/erc-20-indexer/config"
	"github.com/avelex/erc-20-indexer/internal/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	SaveEvent(ctx context.Context, event *abi.Erc20Transfer) error
}

type Indexer struct {
	config config.Config
	repo   Repository
}

func New(config config.Config, repo Repository) *Indexer {
	return &Indexer{
		config: config,
		repo:   repo,
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	client, err := ethclient.DialContext(ctx, i.config.RPC)
	if err != nil {
		return fmt.Errorf("failed to dial to rpc: %w", err)
	}

	defer client.Close()

	erc20, err := abi.NewErc20(i.config.TokenAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create erc20 instance: %w", err)
	}

	startBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			currentBlock, err := client.BlockNumber(ctx)
			if err != nil {
				continue
			}

			if startBlock == currentBlock {
				continue
			}

			opts := &bind.FilterOpts{
				Start: startBlock,
				End:   &currentBlock,
			}

			log.Info().Msgf("Searching for events from block %d to %d", opts.Start, *opts.End)

			iter, err := erc20.FilterTransfer(opts, nil, nil)
			if err != nil {
				continue
			}

			for iter.Next() {
				transfer := iter.Event

				log.Info().Msgf("From: %+v, To: %+v, Value: %+v TxHash: %+v", transfer.From, transfer.To, transfer.Value, transfer.Raw.TxHash)

				if err := i.repo.SaveEvent(ctx, transfer); err != nil {
					continue
				}
			}

			startBlock = currentBlock
		}
	}
}
