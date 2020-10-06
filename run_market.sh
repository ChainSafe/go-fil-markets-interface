#!/usr/bin/env bash

set -o xtrace

LOTUS_DIR=/app/lotus DOCKER= source ./init_market_env.sh
go build -tags 2k -o go-fil-market ./cmd/market
./go-fil-market daemon
