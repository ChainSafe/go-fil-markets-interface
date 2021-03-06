#!/usr/bin/env bash

set -o xtrace

LOTUS_DIR=${LOTUS_DIR:=./extern/lotus}

function cleanup {
    rm -rf ~/.lotus ~/.genesis-sector
    rm localnet.json devgen.car
    # Preserve the lotus node logs
}
trap cleanup EXIT
cleanup

export LOTUS_SKIP_GENESIS_CHECK=_yes_
$LOTUS_DIR/lotus fetch-params 2048
$LOTUS_DIR/lotus-seed pre-seal --sector-size 2KiB --num-sectors 2
$LOTUS_DIR/lotus-seed genesis new localnet.json
$LOTUS_DIR/lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json

if [ -z "$DOCKER" ]
then
    echo "Using default lotus config file"
else
    echo "Using docker lotus config file"
    mkdir ~/.lotus && cp /app/go-fil-markets/config/lotus/config.toml ~/.lotus/config.toml
fi

echo "Starting lotus daemon"
$LOTUS_DIR/lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false 2>&1 | tee lotus.log
