package log

import (
	"fmt"
	"os"
	"time"

	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

//Log for global logging purposes
var Log zerolog.Logger

//InitLog initializes global logging settings
func InitLog() {
	config := util.GetConfig()
	zerolog.TimeFieldFormat = time.RFC3339

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

	Log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	if config.Log.Console {
		Log = Log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	if config.Log.Caller {
		Log = Log.With().Caller().Logger()
	}

	zlog.Logger = Log
}

func init() {
	InitLog()
}
