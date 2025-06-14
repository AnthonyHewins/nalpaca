package conf

import (
	"errors"

	"github.com/nats-io/nats.go"
)

type NATS struct {
	User     string `env:"NATS_USER"`
	Password string `env:"NATS_PASSWORD"`

	URL string `env:"NATS_URL" envDefault:"localhost:4222"`

	Compression bool `env:"NATS_USE_COMPRESSION"`
}

func (b *Bootstrapper) NATSConn(n *NATS) (*nats.Conn, error) {
	l := b.Logger.With(
		"url", n.URL,
		"user", n.User,
		"len(pass)>0", len(n.Password) > 0,
	)

	opts := []nats.Option{nats.Compression(n.Compression)}

	if n.User != "" {
		if n.Password == "" {
			msg := "config error: you set NATS password, but not a user. Please set both, or neither"
			l.Error(msg, "user", n.User)
			return nil, errors.New(msg)
		}

		l.Debug("user credentials set, adding as option")
		opts = append(opts, nats.UserInfo(n.User, n.Password))
	} else if n.Password != "" {
		msg := "config error: you set NATS password, but not a user. Please set both, or neither"
		l.Error(msg)
		return nil, errors.New(msg)
	}

	nc, err := nats.Connect(n.URL, opts...)
	if err != nil {
		l.Error("nats failed connection", "err", err)
		return nil, err
	}

	l.Info("connected to NATS")
	return nc, nil
}
