package util

import (
	"testing"

	"github.com/rs/zerolog/log"
)

func TestLog(t *testing.T) {
	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")

	//these would fail the test
	// Log.Fatal().Msg("fatal")
	// Log.Panic().Msg("panic")
}
