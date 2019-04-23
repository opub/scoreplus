package util

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

//InitLog initializes global logging settings
func InitLog() {
	config := GetConfig()

	//set log level
	if config.Log.Level == "none" || config.Log.Level == "disabled" {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		level, err := zerolog.ParseLevel(config.Log.Level)
		if err != nil {
			fmt.Printf("invalid log level: %s\n", err)
			panic("log initialization failed")
		}
		zerolog.SetGlobalLevel(level)
	}

	//set log format
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true

	//create logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if config.Log.Console {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	if config.Log.Caller {
		logger = logger.With().Caller().Logger()
	}
	log.Logger = logger
}

func init() {
	InitLog()
}
