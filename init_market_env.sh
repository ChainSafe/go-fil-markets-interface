#!/usr/bin/env bash

NODE_PORT=1234
MINER_PORT=2345
MARKET_PORT=8888

ENV_IP=${ENV_IP:=0.0.0.0}

LOTUS_NODE_ADDR=$ENV_IP:$NODE_PORT/rpc/v0
while true; do
    http_code=$(curl --output /dev/null -w ''%{http_code}'' --fail $LOTUS_NODE_ADDR --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "Filecoin.Version", "params": [], "id": 1}')
    if [ "$http_code" -eq 200 ]; then
        echo "Lotus node is up"
        break
    fi
    printf 'Waiting for lotus daemon initialization'
    sleep 5
done

LOTUS_MINER_ADDR=$ENV_IP:$MINER_PORT/rpc/v0
while true; do
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
CONTAINER_NAME=${CONTAINER_NAME:=go-fil-markets-e2e}

NODE_TOKEN=""
if [ -z "$DOCKER" ]; then
    NODE_TOKEN=$($LOTUS_BIN auth create-token --perm admin)
else
    NODE_TOKEN=$(docker exec $CONTAINER_NAME /app/lotus/lotus auth create-token --perm admin | tr -d '\r')
fi

if [[ "$NODE_TOKEN" =~ "ERROR" ]]; then
    echo "$NODE_TOKEN"
    exit 1
fi

NODE_API_INFO=$NODE_TOKEN:/ip4/$ENV_IP/tcp/$NODE_PORT/ws
MARKET_API_INFO=$NODE_TOKEN:/ip4/$ENV_IP/tcp/$MARKET_PORT/ws

echo "NODE_API_INFO=$NODE_API_INFO"
export NODE_API_INFO=$NODE_API_INFO
echo "MARKET_API_INFO=$MARKET_API_INFO"
export MARKET_API_INFO=$MARKET_API_INFO

LOTUS_MINER_BIN=$LOTUS_DIR/lotus-miner

MINER_TOKEN=""
if [ -z "$DOCKER" ]; then
    MINER_TOKEN=$($LOTUS_MINER_BIN auth create-token --perm admin)
else
    MINER_TOKEN=$(docker exec $CONTAINER_NAME /app/lotus/lotus-miner auth create-token --perm admin | tr -d '\r')
fi

if [[ "$MINER_TOKEN" =~ "ERROR" ]]; then
    echo "$MINER_TOKEN"
    exit 1
fi

MINER_API_INFO=$MINER_TOKEN:/ip4/$ENV_IP/tcp/$MINER_PORT/ws
echo "MINER_API_INFO=$MINER_API_INFO"
export MINER_API_INFO=$MINER_API_INFO
