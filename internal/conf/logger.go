package conf

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	Exporter string `env:"LOG_EXPORTER" envDefault:""`
	Level    string `env:"LOG_LEVEL"    envDefault:"info"`
	Fmt      string `env:"LOG_FMT"      envDefault:"json"`
	Src      bool   `env:"LOG_SRC"      envDefault:"true"`
}

func (l Logger) Slog() (*slog.Logger, error) {
	var lvl slog.HandlerOptions
	switch strings.ToLower(l.Level) {
	case "":
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	case "debug":
		lvl = slog.HandlerOptions{Level: slog.LevelDebug, AddSource: l.Src}
	case "info":
		lvl = slog.HandlerOptions{Level: slog.LevelInfo, AddSource: l.Src}
	case "warn":
		lvl = slog.HandlerOptions{Level: slog.LevelWarn, AddSource: l.Src}
	case "err":
		lvl = slog.HandlerOptions{Level: slog.LevelError, AddSource: l.Src}
	default:
		return nil, fmt.Errorf("invalid log level: %s", l.Level)
	}

	exporter, err := l.exporter()
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(l.Fmt) {
	case "", "json":
		return slog.New(slog.NewJSONHandler(exporter, &lvl)), nil
	case "text", "logfmt":
		return slog.New(slog.NewTextHandler(exporter, &lvl)), nil
	}

	return nil, fmt.Errorf("invalid handler format: %s", l.Fmt)
}

func (l Logger) exporter() (io.Writer, error) {
	switch strings.ToLower(l.Exporter) {
	case "":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	}

	file, err := os.Create(l.Exporter)
	if err != nil {
		return nil, err
	}

	return file, nil
}
