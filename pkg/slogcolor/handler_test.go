package slogcolor

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"
)

func Example() {
	opts := DefaultOptions
	opts.Level = slog.LevelDebug
	slog.SetDefault(slog.New(NewHandler(os.Stderr, opts)))

	slog.Info("Initializing")
	slog.Debug("Init done", "duration", 500*time.Millisecond)
	slog.Warn("Slow request!", "method", "GET", "path", "/api/users", "duration", 750*time.Millisecond)
	slog.Error("DB connection lost!", "err", errors.New("connection reset"), "db", "horalky")
	// Output:
}

func BenchmarkLog(b *testing.B) {
	b.StopTimer()
	l := slog.New(NewHandler(os.Stderr, DefaultOptions))

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		l.Info("benchmarking", "i", i)
		b.StopTimer()
	}
}
