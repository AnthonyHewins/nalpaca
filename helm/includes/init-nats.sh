echo Using server ${NATS_URL:?Must set to point to correct NATS location}

if [[ $NATS_USER != "" ]]; then
    echo Using NATS user $NATS_USER
    ${NATS_PASSWORD:?must set password if using NATS_USER, either unset user to null or empty, or supply password}
    exit 1
fi

function create_stream_consumer() {
    if [[ $2 != "true" ]]; then
        echo stream_consumer $1 is marked disabled, got $2. Skipping...
        return
    fi

    echo "==============================================="
    echo "$1"
    echo "==============================================="

    local streamconf=/conf/$1-stream.json
    local consumerconf=/conf/$1-consumer.json

    for i in $streamconf $consumerconf; do
        if [ ! -f $i ]; then
            echo "ERR: config file $i not found"
            exit 1
        fi

        echo "Found config $i..."
    done

    streamname=nalpaca-$1-stream-v0
    echo Creating $streamname with the config below
    cat $streamconf
    nats stream add $streamname --config $streamconf
    echo
    echo

    consumername=nalpaca-$1-consumer-v0 
    echo Creating $consumername under stream $streamname with the config below
    cat $consumerconf
    nats consumer add $streamname $consumername --config $consumerconf
    echo "\n"
}

create_stream_consumer "orders" $ORDERS_ENABLED
create_stream_consumer "tradeupdater" $TRADEUPDATER_ENABLED
create_stream_consumer "cancel" $CANCEL_ENABLED
create_stream_consumer "stockstream" $STOCKSTREAM_ENABLED

if [[ $TRADEUPDATER_ENABLED != "true" ]]; then
    echo "Trade updater disabled, so not creating KV bucket"
    exit 0
fi

echo "=================================================="
echo "KV bucket"
echo "=================================================="
nats kv add ${BUCKET:?Must have a bucket defined to make KV, or disable trade updater}