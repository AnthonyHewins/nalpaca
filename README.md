# Nalpaca

![Nalpaca](./assets/logo.jpg)

- [Nalpaca](#nalpaca)
  - [Deploying](#deploying)
    - [Helm](#helm)
  - [Client usage](#client-usage)
    - [Trades](#trades)
    - [Cancels (in progress)](#cancels-in-progress)
    - [Trade updates by consumer](#trade-updates-by-consumer)
    - [Position updates by KV store](#position-updates-by-kv-store)
  - [Running locally](#running-locally)
    - [Spinning up](#spinning-up)
    - [Debugging](#debugging)

Alpaca over NATS, with k8s support via helm. Nalpaca steps in front of the Alpaca API and makes certain things available to you in NATS.
Nalpaca will take on the onus of trade retries, trade updates, backoffs, failures, metrics, logging and basically everything that you'd probably want
out of your own code interacting with alpaca. It uses protobuf to make messaging as fast and small as possible.
Here are some example things it can do:

- **Trades**: push trades on the NATS message bus, then the trader will perform these trades for you, retrying as much as you want
- **Trade updates**: when a trade is filled, it will push notifications
- **Positions**: positions are available in a KV store
- **Trade cancels**: cancel trades (in progress)
- **Streaming**: stream stocks, and eventually options

## Deploying

Deployment is geared toward k8s but not required. Right now, this only works having a single pod, so it can't handle any degree of scaling yet if you use the `tradeupdater` component because each pod would duplicate trade updates. However, cancels and order creation can be scaled

### Helm

The only real dependency for nalpaca is that you have an account and some NATS jetstream client to connect to. If you have those 2 things, you're good to go

1. Set the NATS url in `nats.url` you want to connect to
2. Set the api key `alpaca.apiKey`
3. Set the `secrets.name` value (default: `nalpaca-secrets`) to a k8s secret containing the key `ALPACA_API_SECRET`
4. Optional: set a subject prefix so all subjects are namespaced for nalpaca (default: `nalpaca`)
5. Optional: if you want to go right into live trading, uncomment the production API URL

**Adding user/password:**

1. Set the username in the NATS config block `nats.user`
2. Add `NATS_PASSWORD` to the secret located at `secrets.name`

## Client usage

The client only will need access to a NATS jetstream connection, and NATS KV client if you need any of the KV features

### Trades

To create trades, create a protobuf message of type [`tradesvc.v0.Trade`](./api/proto/tradesvc/v0/tradesvc.proto). Then create an idempotency ID that's <= 128 chars (a limit set by alpaca). Then send it on the topic `<prefix>.orders.v0.<string client order id (<=128 chars)>`, where the prefix is the prefix set in the deployment (default: `nalpaca`)

### Cancels (in progress)

To execute a cancel, you just publish an empty message on `<prefix>.cancel.v0.<order ID or special keyword "ALL">`. Using special keyword `ALL` will initiate a cancel of all orders

### Trade updates by consumer

Updates on trades are received as a consumer. Connect to `nalpaca-tradeupdater-consumer-v0`. Once connected, the consumer will receive updates on subject `<prefix>.tradeupdates.v0.<TICKER>` with message type [`tradesvc.v0.TradeUpdate`](./api/proto/tradesvc/v0/tradesvc.proto)

### Position updates by KV store

Positions for the account are placed in a KV store under the bucket specified in the helm deployment (default bucket name: `nalpaca`). Once you connect to the bucket the key `positions` will store the positions of the account with type [`tradesvc.v0.Positions`](./api/proto/tradesvc/v0/tradesvc.proto). Being that it's a KV store you can subscribe to changes or choose to fetch whenever desired

## Running locally

**Make tasks**

```
(target)                       Build a target binary in current arch for running locally: nalpaca
all                            Build all targets
docker                         build docker image
compose                        build docker compose
run-compose                    Run a binary with docker compose
proto                          buf generate
clean                          gofmt, go generate, then go mod tidy, and finally rm -rf bin/
test                           Run go vet, then test all files
help                           Print help
```

### Spinning up

1. `cp .env.tmpl .env`
2. Set all your secrets in `.env` (nats user/pass not required)
3. Choose the method:

**Docker compose**

1. Optional: if you feel the need to use a different NATS connection, you'll want to edit the `compose.yaml` to do so
2. `make compose run-compose`

**Terminal**

This option requires NATS running elsewhere that you can connect to, and requires running the init script

1. Run the `docker/nats-init.sh` script; you can edit the options to your liking for testing
2. `make all` and `make run-nalpaca` will start the server

### Debugging

Debugging PBF requires some scripting. Some scripts exist in the `scripts` directory to debug certain things, e.g. debugging the positions KV store