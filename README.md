# Nalpaca

![Nalpaca](./assets/logo.jpg)

Alpaca with NATS. Nalpaca steps in front of the Alpaca API and makes certain things available to you in NATS.
The goal is to not directly use the SDK, handle retries, backoffs, failures, and all the rest in your code but
rather delegate that responsibility to nalpaca, which will make these events available on your configured NATS
connection

- **Trades**: push trades on the NATS message bus, then the trader will perform these trades for you, retrying as much as you want
- **Trade updates**: when a trade is filled, it will push notifications

## API

Messages to and from NATS are all protobuf. See [`api/proto`](./api/proto/) for the protobuf definitions

## Usage

There are several binaries that are built that compose a nalpaca deployment that you can use.
Subjects are prefixed by a configurable prefix that defaults to `nalpaca`, and the subjects are broken down as follows:

### orders

Can handle order creation and cancels

**Push** to these subjects:

- `<prefix>.orders.v0.<string client order id (<=128 chars)>`: push [`tradesvc.v0.Trade`](./api/proto/tradesvc/v0/tradesvc.proto) on this subject to perform a trade on the logged in portfolio. You'll need to generate a client ID first
- `<prefix>.cancel.v0.<order ID or special keyword "ALL">`: push on this subject to cancel the specific order (no payload required). Using special keyword `ALL` will initiate a cancel of all orders

### tradeupdater

Updates on trades

**Connect as a consumer**:

- `nalpaca-tradeupdater-consumer-v0` (updates on `<prefix>.tradeupdates.v0.<TICKER>`): subscribe to this subject range to get trade updates for the specified ticker specified by [`tradesvc.v0.TradeUpdate`](./api/proto/tradesvc/v0/tradesvc.proto)

## Running

**Docker compose**

This is the easiest way with no dependencies, minus Go/docker compose: `make compose run-compose` from the repo root will build the project and run it with docker compose. After building you can skip re-building with `make run-compose`

**Terminal**

TODO

1. Have your Alpaca keys ready and available
2. Install Go, NATS server
3. Make sure that NATS server is running
4. (add something about initializing nats)...
5. `make all` and `./bin/<desired binary>` will get you up and running