package log

import "testing"

func TestLog(t *testing.T) {

	Log.Debug().Msg("debug")
	Log.Info().Msg("info")
	Log.Warn().Msg("warn")
	Log.Error().Msg("error")

	//these would fail the test
	// Log.Fatal().Msg("fatal")
	// Log.Panic().Msg("panic")

}
