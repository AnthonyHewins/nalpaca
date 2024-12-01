package conf

import "github.com/nats-io/nats.go"

type NATS struct {
	Disable bool `env:"DISABLE_NATS" envDefault:"false"`

	User     string `env:"NATS_USER"`
	Password string `env:"NATS_PASSOWRD"`

	URL string `env:"NATS_URL" envDefault:"localhost:4222"`

	Compression bool `env:"NATS_USE_COMPRESSION"`
}

func (b *Bootstrapper) NATSConn(n *NATS) (*nats.Conn, error) {
	if n.Disable {
		b.Logger.Info("nats set to be disabled")
		return nil, nil
	}

	l := b.Logger.With(
		"url", n.URL,
		"user", n.User,
		"len(pass)>0", len(n.Password) > 0,
	)

	opts := []nats.Option{nats.Compression(n.Compression)}

	if n.User != "" && n.Password != "" {
		l.Debug("user credentials set, adding as option")
		opts = append(opts, nats.UserInfo(n.User, n.Password))
	}

	nc, err := nats.Connect(n.URL, opts...)
	if err != nil {
		l.Error("nats failed connection", "err", err)
		return nil, err
	}

	return nc, nil
}
