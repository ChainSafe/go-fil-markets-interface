#!/usr/bin/env bash

set -o xtrace

if [ -z "$LOTUS_DIR" ]
then
    LOTUS_DIR=./extern/lotus
fi

function cleanup {
    rm -rf ~/.lotus ~/.genesis-sector
    rm localnet.json devgen.car
    rm -rf "$LOTUS_PATH"
    # Preserve the lotus node logs
}
trap cleanup EXIT
cleanup

if [ -z "$DOCKER" ]
then
    echo "Using default lotus config file"
else
    echo "Using docker lotus config file"
    mkdir -p "$LOTUS_PATH" && cp /app/go-fil-markets/config/lotus/config.toml "$LOTUS_PATH"/config.toml
fi

export LOTUS_SKIP_GENESIS_CHECK=_yes_
$LOTUS_DIR/lotus fetch-params 2048
$LOTUS_DIR/lotus-seed pre-seal --sector-size 2KiB --num-sectors 2
$LOTUS_DIR/lotus-seed genesis new localnet.json
$LOTUS_DIR/lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json

echo "Starting lotus daemon"
$LOTUS_DIR/lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false --config="$LOTUS_PATH"/config.toml 2>&1 | tee lotus.log
