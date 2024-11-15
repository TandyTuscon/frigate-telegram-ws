#!/bin/sh

CONFIG_PATH="/root/config/config.yml"

if [ ! -f "$CONFIG_PATH" ]; then
    echo "Default config.yml not found. Creating..."
    mkdir -p /root/config
    cp /root/default_config.yml $CONFIG_PATH
fi

exec "$@"
