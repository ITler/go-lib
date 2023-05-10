package targets

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/itler/go-lib/magefiles/deps"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Ci installs dependencies in a quick way, suitable for temporary pipeline runners
func Ci(ctx context.Context) error {
	success, err := deps.InstallDependencies(ctx)
	if !success {
		return err
	}
	if err != nil {
		log.Warn().Msg(fmt.Errorf("install dependencies incomplete: %w", err).Error())
	}
	return nil
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if err := deps.CheckDependencies(context.Background()); err != nil {
		log.Warn().Msg(fmt.Errorf("Dependency not installed - %w", err).Error())
	}
}
