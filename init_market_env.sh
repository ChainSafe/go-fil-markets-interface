#!/usr/bin/env bash

NODE_PORT=1234
MINER_PORT=2345
MARKET_PORT=8888

LOTUS_NODE_ADDR=127.0.0.1:$NODE_PORT/rpc/v0
while true
do
    http_code=$(curl --output /dev/null -w ''%{http_code}'' --fail $LOTUS_NODE_ADDR --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "Filecoin.Version", "params": [], "id": 1}')
    if [ "$http_code" -eq 200 ]; then
		  echo "Lotus node is up"
		  break
	  fi
    printf 'Waiting for lotus daemon initialization'
    sleep 5
done

LOTUS_MINER_ADDR=127.0.0.1:$MINER_PORT/rpc/v0
while true
do
    http_code=$(curl --output /dev/null -w ''%{http_code}'' --fail $LOTUS_MINER_ADDR --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "Filecoin.Version", "params": [], "id": 1}')
    if [ "$http_code" -eq 200 ]; then
		  echo "Lotus miner is up"
		  break
	  fi
    printf 'Waiting for lotus miner initialization'
    sleep 5
done

LOTUS_DIR=${LOTUS_DIR:=./extern/lotus}

LOTUS_BIN=$LOTUS_DIR/lotus
NODE_TOKEN=$($LOTUS_BIN auth create-token --perm admin)

if [[ "$NODE_TOKEN" =~ "ERROR" ]]; then
    echo $NODE_TOKEN
    exit 1
fi

LOTUS_MINER_BIN=$LOTUS_DIR/lotus-miner
MINER_TOKEN=$($LOTUS_MINER_BIN auth create-token --perm admin)

if [[ "$MINER_TOKEN" =~ "ERROR" ]]; then
    echo $MINER_TOKEN
    exit 1
fi

export NODE_API_INFO=$NODE_TOKEN:/ip4/127.0.0.1/tcp/$NODE_PORT/ws
export MARKET_API_INFO=$NODE_TOKEN:/ip4/127.0.0.1/tcp/$MARKET_PORT/ws
export MINER_API_INFO=$MINER_TOKEN:/ip4/127.0.0.1/tcp/$MINER_PORT/ws
