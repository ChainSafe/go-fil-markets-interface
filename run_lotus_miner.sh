#!/usr/bin/env bash

set -o xtrace

function cleanup {
    rm -rf ~/.lotusminer
    # Preserve the lotus miner logs
}
trap cleanup EXIT
cleanup

LOTUS_DIR=./extern/lotus
LOTUS_NODE_URL=127.0.0.1:1234/rpc/v0
while true
do
    http_code=$(curl --output /dev/null -w ''%{http_code}'' --fail $LOTUS_NODE_URL --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "Filecoin.Version", "params": [], "id": 1}')
    if [ "$http_code" -eq 200 ]; then
		  echo "Lotus node is up"
		  break
	  fi
    printf 'Waiting for lotus daemon initialization'
    sleep 5
done
echo "Initializing lotus wallet"
$LOTUS_DIR/lotus wallet import --as-default ~/.genesis-sectors/pre-seal-t01000.key
echo "Initializing lotus miner"
$LOTUS_DIR/lotus-miner init --genesis-miner --actor=t01000 --sector-size=2KiB --pre-sealed-sectors=~/.genesis-sectors --pre-sealed-metadata=~/.genesis-sectors/pre-seal-t01000.json --nosync
echo "Starting lotus miner"
$LOTUS_DIR/lotus-miner run --nosync 2>&1 | tee lotus_miner.log
