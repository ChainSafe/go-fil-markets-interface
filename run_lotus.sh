#!/bin/sh

set -o xtrace
rm -rf ~/.lotus ~/.genesis-sectors
export LOTUS_SKIP_GENESIS_CHECK=_yes_
LOTUS_DIR=./extern/lotus
$LOTUS_DIR/lotus fetch-params 2048
$LOTUS_DIR/lotus-seed pre-seal --sector-size 2KiB --num-sectors 2
$LOTUS_DIR/lotus-seed genesis new localnet.json
$LOTUS_DIR/lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json
echo "Starting lotus daemon"
$LOTUS_DIR/lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false 2>&1 | tee lotus.log
