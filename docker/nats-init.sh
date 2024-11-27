echo "Creating NATS config"

set -x

SERVER="${SERVER:-nats://nats:4222}"
SUBJECT_PREFIX="${SUBJECT_PREFIX:-nalpaca}"

STORAGE="${STORAGE:memory}"

# Order stream
ORDER_STREAM_MAX_AGE="${ORDER_STREAM_MAX_AGE:-2m}"
ORDER_STREAM_MAX_BYTES="${ORDER_STREAM_MAX_BYTES:-10240}"
ORDER_STREAM_BACKOFF="${ORDER_STREAM_MAX_BYTES:-1s,3s,1m}"
ORDER_STREAM_MAX_DELIVER="${ORDER_STREAM_MAX_DELIVER:-3}"
ORDER_STREAM_MAX_MSGS="${ORDER_STREAM_MAX_MSGS:100}"
ORDER_STREAM_MAX_MSG_SIZE="${ORDER_STREAM_MAX_MSG_SIZE:-512}"
ORDER_STREAM_MAX_CONSUMERS="${ORDER_STREAM_MAX_CONSUMERS:--1}"

set +x

ORDER_STREAM=nalpaca-order-stream-v0 
SUBJECT="$SUBJECT_PREFIX.orders.v0.*"
nats stream add $ORDER_STREAM \
      --subjects $SUBJECT \
      --description="Nalpaca order stream" \
      --storage=$STORAGE \
      --compression=s2 \
      --dupe-window="20s" \
      --ack \
      --no-allow-rollup \
      --retention=work \
      --discard=old \
      --replicas $ORDER_STREAM_REPLICAS \
      --max-age $ORDER_STREAM_MAX_AGE \
      --max-bytes $ORDER_STREAM_MAX_BYTES \
      --max-consumers $ORDER_STREAM_MAX_CONSUMERS \
      --max-msg-size $ORDER_STREAM_MAX_MSG_SIZE \
      --max-msgs $ORDER_STREAM_MAX_MSGS

nats consumer add $ORDER_STREAM nalpaca-order-consumer-v0 \
    --description "Nalpaca order stream consumer. Nalpaca will consume messages on this stream and send them." \
    --max-deliver $ORDER_STREAM_MAX_DELIVER \
    --filter $SUBJECT \
    --backoff-steps $ORDER_STREAM_BACKOFF