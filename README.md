# Nalpaca

![Nalpaca](./assets/logo.jpg)

Alpaca with NATS. Nalpaca steps in front of the Alpaca API and makes certain things available to you in NATS:

- **Trades**: push trades on the NATS message bus, then the trader will perform these trades for you, retrying as much as you want
- **Trade updates**: when a trade is filled, it will push notifications and update your portfolio in a KV store

## API

Messages to and from NATS are all protobuf. See [`api/proto`](./api/proto/) for the protobuf definitions

## Subjects

Subjects are prefixed by a configurable prefix that defaults to `nalpaca`, and the subjects are broken down as follows:

**Push** to these subjects:

- `<prefix>.orders.v0.<string client order id (<=128 chars)>`: push [`tradesvc.v0.Trade`](./api/proto/tradesvc/v0/tradesvc.proto) on this subject to perform a trade on the logged in portfolio. You'll need to generate a client ID first
- `<prefix>.cancel.v0.<order ID or special keyword "ALL">`: push on this subject to cancel the specific order. Using special keyword `ALL` will initiate a cancel of all orders

**Pull** to these subjects:

- `<prefix>.tradeupdates.v0.<TICKER>`: subscribe to this subject range to get trade updates for the specified ticker specified by [`tradesvc.v0.TradeUpdate`](./api/proto/tradesvc/v0/tradesvc.proto)

## Running

**Docker compose**

This is the easiest way with no dependencies, minus Go/docker compose: `make compose run-compose` from the repo root will build the project and run it with docker compose. After building you can skip re-building with `make run-compose`

**Terminal**

You'll need NATS installed and running with Golang