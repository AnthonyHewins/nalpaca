package conf

import "log/slog"

type Bootstrapper struct {
	Logger *slog.Logger
}

func NewBootstrapper(l *Logger) (*Bootstrapper, error) {
	logger, err := l.Slog()
	if err != nil {
		return nil, err
	}

	return &Bootstrapper{Logger: logger}, nil
}
