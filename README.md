# Nalpaca

![Nalpaca](./assets/logo.jpg)

- [Nalpaca](#nalpaca)
  - [Components](#components)
    - [Trades](#trades)
    - [Cancels (in progress)](#cancels-in-progress)
    - [Trade updates](#trade-updates)
    - [Streaming bars](#streaming-bars)
      - [Option 1: receive as consumer](#option-1-receive-as-consumer)
      - [Option 2: KV store](#option-2-kv-store)
  - [Deploying](#deploying)
    - [Helm](#helm)
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

## Components

Nalpaca is based on NATS resources, there are 3 primary streams to look at that it creates and a myriad of consumers, and a single KV store

1. **Action stream** (`nalpaca-action-stream`, subjects `nalpaca.action.>`): this stream is the stream that you will write to if you want to have nalpaca perform some actions on your behalf like executing trades. View the docs for the component you wish to use to see what subjects to publish to in order to use it
  1. The action consumer is also created, which is configured as a `workqueue` and is reserved only for nalpaca to use
1. **Data stream** (`nalpaca-data-stream`, subjects `nalpaca.data.>`): public market data that nalpaca receives from alpaca and republishes — stock bars, quotes and trades, option quotes and trades, and news. Nalpaca writes to this stream, and any number of consumers it offers will be what you subscribe to in order to get info on it. This stream is configured with `interest` retention
1. **Account stream** (`nalpaca-account-stream`, subjects `nalpaca.account.>`): events about *your* account rather than the market, such as fills on orders you placed. Kept separate from the data stream so that market data subscribers don't receive your private account activity, and because these events span every asset class rather than just stocks. Also configured with `interest` retention
1. **KV store**: this KV store is a global store for the app, where it will write useful things like your current position list, which you can subscribe to via keywatcher, or just access on the fly if needed. See the docs for the component you want to use for more

Because both data-bearing streams use `interest` retention, a message is only kept if some consumer's filter matches it. If you publish to a subject no consumer is listening on, the message is acknowledged and then dropped — so add a consumer before expecting data to accumulate.

### Subjects

| Subject | Contents |
| --- | --- |
| `nalpaca.action.orders.<idempotency key>` | orders you want nalpaca to place |
| `nalpaca.action.cancel.<order ID or ALL>` | cancels you want nalpaca to execute |
| `nalpaca.data.stocks.bars.<TICKER>` | stock bars |
| `nalpaca.data.stocks.quotes.<TICKER>` | stock quotes |
| `nalpaca.data.stocks.trades.<TICKER>` | stock trades (market-wide, not yours) |
| `nalpaca.data.options.quotes.<CONTRACT>` | option quotes |
| `nalpaca.data.options.trades.<CONTRACT>` | option trades |
| `nalpaca.data.news.<ID>` | news articles |
| `nalpaca.account.tradeupdates.<TICKER>` | fills and other updates on *your* orders |

### Trades

To create trades, create a protobuf message of type [`tradesvc.v0.Trade`](./api/proto/tradesvc/v0/tradesvc.proto). Then create an idempotency ID and send it on the topic `nalpaca.action.orders.<string client order id (<=128 chars)>`

### Cancels (in progress)

To execute a cancel, you just publish an empty message on `<prefix>.action.cancel.<order ID or special keyword "ALL">`. Using special keyword `ALL` will initiate a cancel of all orders

### Trade updates

TODO docs

### Streaming bars

Stream bar data. Possibly the most useful feature, allows you to broadcast messages across your architecture for lots of listeners

#### Option 1: receive as consumer

Updates on trades can be received as a consumer. Connect to `nalpaca-tradeupdate-consumer`. Once connected, the consumer will receive updates on subject `<prefix>.account.tradeupdates.<TICKER>` with message type [`tradesvc.v0.TradeUpdate`](./api/proto/tradesvc/v0/tradesvc.proto)

#### Option 2: KV store

Positions for the account are placed in a KV store under the bucket specified in the helm deployment (default bucket name: `nalpaca`). Once you connect to the bucket the key `positions` will store the positions of the account with type [`tradesvc.v0.Positions`](./api/proto/tradesvc/v0/tradesvc.proto). Being that it's a KV store you can subscribe to changes or choose to fetch whenever desired

## Deploying

Deployment is geared toward k8s but not required. Right now, this only works having a single pod, so it can't handle any degree of scaling yet because of the nature of alpaca websockets. However, future deployments may be able to

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

## Running locally

**Make tasks**

```
(target)                       Build a target binary in current arch for running locally: nalpaca
all                            Build all targets
docker                         build docker image 
compose                        build docker compose
up                             run docker compose
down                           teardown docker compose
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
2. `make compose up`

**Terminal**

This option requires NATS running elsewhere that you can connect to, and requires running the init script

1. Run the [`scripts/nats/init.sh`](scripts/nats/init.sh) script; you can edit the options to your liking for testing
2. `make all run-nalpaca` will start the server

### Debugging

Debugging PBF requires some scripting. Some scripts exist in the `scripts` directory to debug certain things, e.g. debugging the positions KV store