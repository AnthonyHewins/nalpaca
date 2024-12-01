echo "Creating NATS config"

set -x

SERVER="${SERVER:-nats://nats:4222}"
SUBJECT_PREFIX="${SUBJECT_PREFIX:-nalpaca}"

STORAGE="${STORAGE:-memory}"

# Order stream
ORDER_STREAM_MAX_AGE="${ORDER_STREAM_MAX_AGE:-2m}"
ORDER_STREAM_MAX_BYTES="${ORDER_STREAM_MAX_BYTES:-10240}"
ORDER_STREAM_BACKOFF="${ORDER_STREAM_MAX_BYTES:-1s,3s,1m}"
ORDER_STREAM_MAX_DELIVER="${ORDER_STREAM_MAX_DELIVER:-3}"
ORDER_STREAM_MAX_MSGS="${ORDER_STREAM_MAX_MSGS:-100}"
ORDER_STREAM_REPLICAS="${ORDER_STREAM_REPLICAS:-1}"
ORDER_STREAM_MAX_MSG_SIZE="${ORDER_STREAM_MAX_MSG_SIZE:-512}"
ORDER_STREAM_MAX_CONSUMERS="${ORDER_STREAM_MAX_CONSUMERS:--1}"

# Trade updater
TRADE_UPDATER_STREAM_MAX_AGE="${TRADE_UPDATER_STREAM_MAX_AGE:-1h}"
TRADE_UPDATER_STREAM_MAX_BYTES="${TRADE_UPDATER_STREAM_MAX_BYTES:-51200}"
TRADE_UPDATER_STREAM_BACKOFF="${TRADE_UPDATER_STREAM_MAX_BYTES:-1s,3s,5s,10s}"
TRADE_UPDATER_STREAM_MAX_DELIVER="${TRADE_UPDATER_STREAM_MAX_DELIVER:-4}"
TRADE_UPDATER_STREAM_MAX_MSGS="${TRADE_UPDATER_STREAM_MAX_MSGS:-100}"
TRADE_UPDATER_STREAM_MAX_MSG_SIZE="${TRADE_UPDATER_STREAM_MAX_MSG_SIZE:-512}"

#=================================
# Orders
#=================================
STREAM=nalpaca-order-stream-v0 
SUBJECT="$SUBJECT_PREFIX.orders.v0.*"
echo Creating stream $STREAM
nats stream add $STREAM \
      -j \
      --subjects="$SUBJECT" \
      --description="Nalpaca order stream" \
      --storage=$STORAGE \
      --compression=s2 \
      --dupe-window="20s" \
      --no-allow-rollup \
      --retention=work \
      --discard=old \
      --replicas="$ORDER_STREAM_REPLICAS" \
      --max-age $ORDER_STREAM_MAX_AGE \
      --max-bytes $ORDER_STREAM_MAX_BYTES \
      --max-consumers="$ORDER_STREAM_MAX_CONSUMERS" \
      --max-msgs-per-subject=0 \
      --max-msg-size $ORDER_STREAM_MAX_MSG_SIZE \
      --max-msgs $ORDER_STREAM_MAX_MSGS

exit 1

CONSUMER=nalpaca-order-consumer-v0
echo Creating consumer $CONSUMER
nats consumer add $STREAM $CONSUMER \
    --description "Nalpaca order stream consumer. Nalpaca will consume messages on this stream and send them." \
    --max-deliver $ORDER_STREAM_MAX_DELIVER \
    --filter $SUBJECT \
      -j \
      --replicas 1 \
    --backoff-steps $ORDER_STREAM_BACKOFF

#=================================
# Trade updater
#=================================
STREAM=nalpaca-tradeupdater-stream-v0 
SUBJECT="$SUBJECT_PREFIX.tradeupdates.v0.*"
echo Creating stream $STREAM
nats stream add $STREAM \
      --subjects $SUBJECT \
      --description="Nalpaca trade update stream" \
      --storage=$STORAGE \
      --compression=s2 \
      --dupe-window="20s" \
      --ack \
      --no-allow-rollup \
      --retention=work \
      --replicas 1 \
      --discard=old \
      --replicas $ORDER_STREAM_REPLICAS \
      --max-age $ORDER_STREAM_MAX_AGE \
      --max-bytes $ORDER_STREAM_MAX_BYTES \
      --max-msg-size $ORDER_STREAM_MAX_MSG_SIZE \
      -j \
      --max-msgs $ORDER_STREAM_MAX_MSGS

CONSUMER=nalpaca-tradeupdate-consumer-v0
echo Creating consumer $CONSUMER
nats consumer add $STREAM $CONSUMER \
    --description "Nalpaca order stream consumer. Nalpaca will consume messages on this stream and send them." \
    --max-deliver $ORDER_STREAM_MAX_DELIVER \
      --replicas 1 \
    --filter $SUBJECT \
      -j \
    --backoff-steps $ORDER_STREAM_BACKOFF