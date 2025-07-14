package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(env string) zerolog.Logger {
	var output zerolog.ConsoleWriter
	if env == "dev" {
		output = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		return zerolog.New(output).With().Timestamp().Logger()
	}
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
