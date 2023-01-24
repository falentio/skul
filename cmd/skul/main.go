package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/falentio/skul/internal/app"
	"github.com/falentio/skul/internal/pkg/response"
)

func main() {
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.DebugLevel)

	response.SetLogger(logger)

	opts := app.AppOptions{
		Logger: logger,
	}
	if err := opts.Init(); err != nil {
		logger.Fatal().Err(err).Msg("failed to init options")
	}
	logger.Info().Interface("options", opts).Msg("configured options")

	app := app.Application{
		Options: opts,
		Logger:  logger,
	}

	app.InitRepository()

	logger.Info().Msg("seeding repository")
	if err := app.SeedRepository(); err != nil {
		logger.Fatal().Err(err).Msg("failed while seeding repository")
	}

	logger.Info().Msg("initializing handler")
	app.InitHandler()

	logger.Info().Str("address", opts.Addr).Msg("app server started")
	if err := app.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("failed to run http server")
	}
}
