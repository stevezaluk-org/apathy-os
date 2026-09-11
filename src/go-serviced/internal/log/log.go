package log

import (
	"log/slog"
	"os"
)

// Init Initializes the global slog logger. Extraordinarily simple
//
// TODO: Configuration for this needs to be exposed at config.Config
func Init() {
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	))

	slog.SetLogLoggerLevel(slog.LevelInfo)
}
