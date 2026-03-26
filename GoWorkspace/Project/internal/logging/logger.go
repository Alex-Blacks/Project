package logging

import (
	"log/slog"
	"os"
)

// NewLogger создаёт новый slog.Logger с JSON-хендлером для вывода в stdout
func NewLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
}
