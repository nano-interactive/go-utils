package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Default datetime format
const DateTimeFormat = "2006-01-02 15:04:05"

// Set default options for Zerolog with log level and pretty print
// If pretty print is false, it will write to Standard Output
func ConfigureDefaultLogger(level string, prettyPrint bool, output ...io.Writer) {
	zerologLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("Failed to parse logging level: " + level)
	}

	zerolog.SetGlobalLevel(zerologLevel)
	zerolog.TimeFieldFormat = DateTimeFormat
	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.TimestampFunc = time.Now().UTC

	var w io.Writer

	if prettyPrint {
		w = zerolog.NewConsoleWriter()
	} else {
		w = os.Stdout
	}

	if len(output) > 0 {
		w = zerolog.MultiLevelWriter(w, output[0])
	}

	log.Logger = log.Output(w)
}

// Returns an instance of zerolog.Logger with configured log level and pretty print flag
// If pretty print is false, it will write to Standard Output
func New(level string, prettyPrint bool) zerolog.Logger {
	var logger zerolog.Logger

	zerologLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("Failed to parse logging level: " + level)
	}

	var w io.Writer

	if prettyPrint {
		w = zerolog.NewConsoleWriter()
	} else {
		w = os.Stdout
	}

	logger = zerolog.New(w).
		With().
		Timestamp().
		Logger().
		Level(zerologLevel)

	return logger
}
