${NATS_URL:?Must set to point to correct NATS location}
${CONFIG_PATH:?Must set config path}

if [[ $NATS_USER != "" ]]; then
    ${NATS_PASSWORD:?must set password if using NATS_USER}
fi

echo Using server $NATS_URL
echo Using NATS user $NATS_USER

function create_component() {
    if [[ $2 != "true" ]]; then
        echo Component $1 is marked disabled, got $2. Skipping...
        return
    fi

    streamconf=$CONFIG_PATH/$1-stream.json
    consumerconf=$CONFIG_PATH/$1-consumer.json

    for i in $streamconf $consumerconf; do
        if [ ! -f $i ]; then
            echo "ERR: config file $i not found"
            exit 1
        fi

        echo "Found config $i..."
    done

    echo Creating $1 component...
    set -x
    nats stream add nalpaca-$1-stream-v0 --config $streamconf
    nats consumer add nalpaca-$1-consumer-v0 --config $consumerconf
    set +x
}

create_component "orders" $ORDERS_ENABLED
create_component "tradeupdater" $TRADEUPDATER_ENABLED
create_component "cancel" $CANCEL_ENABLED