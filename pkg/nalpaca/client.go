package nalpaca

import (
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	prefix string
	nc     jetstream.JetStream
}

func NewClient(nc jetstream.JetStream, prefix string) *Client {
	return &Client{nc: nc, prefix: prefix}
}
