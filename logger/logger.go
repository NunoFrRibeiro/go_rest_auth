package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func init() {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
}

func Info(message string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now().UTC().Local(), slog.LevelInfo, fmt.Sprintf(message, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func Debug(message string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now().UTC().Local(), slog.LevelDebug, fmt.Sprintf(message, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func Error(message string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now().UTC().Local(), slog.LevelError, fmt.Sprintf(message, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}
