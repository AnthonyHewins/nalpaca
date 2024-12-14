package nalpaca

import (
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	prefix string
	nc     jetstream.JetStream
	kv     jetstream.KeyValue
}

func NewClient(nc jetstream.JetStream, kv jetstream.KeyValue, prefix string) *Client {
	return &Client{
		prefix: prefix,
		nc:     nc,
		kv:     kv,
	}
}
