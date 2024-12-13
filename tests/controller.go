package main

import (
	"time"

	"github.com/AnthonyHewins/nalpaca/pkg/nalpaca"
	"github.com/nats-io/nats.go"
)

type controller struct {
	timeout time.Duration
	nc      *nats.Conn
	client  *nalpaca.Client
}
