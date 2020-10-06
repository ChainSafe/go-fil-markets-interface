#!/usr/bin/env bash

CONTAINER_NAME=go-fil-markets-e2e
function cleanup {
    # Stop the docker container
    docker stop $CONTAINER_NAME || true
    docker container rm $CONTAINER_NAME || true
}
trap cleanup EXIT
cleanup

# Pull and run the docker image
docker pull arijitad/go-fil-markets:latest
docker run -d --name $CONTAINER_NAME -p 1234:1234 -p 2345:2345 -p 8888:8888 --rm -it arijitad/go-fil-markets

# Initialize the env variables
DOCKER=_yes_ CONTAINER_NAME=$CONTAINER_NAME source ./init_market_env.sh

# Wait for market initialization to complete
ENV_IP=${ENV_IP:=0.0.0.0}
MARKET_PORT=8888
MARKET_ADDR=$ENV_IP:$MARKET_PORT/rpc/v0
while true; do
    http_code=$(curl --output /dev/null -w ''%{http_code}'' --fail $MARKET_ADDR --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "Market.ClientListDataTransfers", "params": [], "id": 1}')
    if [ "$http_code" -eq 200 ]; then
        echo "Market is up"
        break
    fi
    printf 'Waiting for Market daemon initialization'
    sleep 5
done

make submodule

# Run the test
cd cmd/client
go test -v -timeout 60m -run ^TestMarketStorage$
sleep 10s
go test -v -timeout 60m -run ^TestMarketRetrieval$

cleanup
