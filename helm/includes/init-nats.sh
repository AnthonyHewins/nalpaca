echo Using server ${NATS_URL:?Must set to point to correct NATS location}

if [[ $NATS_USER != "" ]]; then
    echo Using NATS user $NATS_USER
    ${NATS_PASSWORD:?must set password if using NATS_USER, either unset user to null or empty, or supply password}
    exit 1
fi

function create_component() {
    if [[ $2 != "true" ]]; then
        echo Component $1 is marked disabled, got $2. Skipping...
        return
    fi

    streamconf=/conf/$1-stream.json
    consumerconf=/conf/$1-consumer.json

    for i in $streamconf $consumerconf; do
        if [ ! -f $i ]; then
            echo "ERR: config file $i not found"
            exit 1
        fi

        echo "Found config $i..."
    done

    echo Creating $1 stream with the config below
    cat $streamconf
    nats stream add nalpaca-$1-stream-v0 --config $streamconf

    echo Creating $1 consumer with the config below
    cat $consumerconf
    nats consumer add nalpaca-$1-consumer-v0 --config $consumerconf
}

create_component "orders" $ORDERS_ENABLED
create_component "tradeupdater" $TRADEUPDATER_ENABLED
create_component "cancel" $CANCEL_ENABLED