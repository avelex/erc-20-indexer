package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/avelex/erc-20-indexer/config"
	"github.com/avelex/erc-20-indexer/internal/indexer"
	"github.com/avelex/erc-20-indexer/internal/repository/timescale"
	"github.com/avelex/erc-20-indexer/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	context.AfterFunc(ctx, func() {
		log.Info().Msg("shutting down indexer")
	})

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	log.Info().Msgf("config: %+v", cfg)

	db, err := pgxpool.New(ctx, cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	defer db.Close()

	repo := timescale.NewRepository(db)
	indexer := indexer.New(cfg, repo)
	handler := server.New(repo)

	go func() {
		http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
			events, err := handler.LastEvents(r.Context(), 10)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(events)
		})

		log.Info().Msg("starting server on port 8080")
		http.ListenAndServe(":8080", nil)
	}()

	log.Info().Msg("starting indexer")
	if err := indexer.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start indexer")
	}
	log.Info().Msg("indexer stopped!")
}
