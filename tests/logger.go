package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type logger struct {
	quiet bool
	level string
	fmt   string
	src   bool
}

func (l logger) slog() (*slog.Logger, error) {
	h, err := l.handler()
	if err != nil {
		return nil, err
	}

	return slog.New(h), nil
}

func (l logger) handler() (slog.Handler, error) {
	lvl, err := l.getLevel()
	if err != nil {
		return nil, err
	}

	opts := slog.HandlerOptions{
		AddSource: l.src,
		Level:     lvl,
	}

	switch strings.ToLower(l.fmt) {
	case "json":
		return slog.NewJSONHandler(os.Stdout, &opts), nil
	case "text", "logfmt":
		return slog.NewTextHandler(os.Stdout, &opts), nil
	default:
		return nil, fmt.Errorf("invalid fmt: %s", l.fmt)
	}
}

func (l logger) getLevel() (slog.Level, error) {
	switch strings.ToLower(l.level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid level %s", l.level)
	}
}
