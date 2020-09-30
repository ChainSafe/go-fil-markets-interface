#!/usr/bin/env bash

make submodule

cp /app/go-fil-markets/config/lotus/config.toml ~/lotus/config.toml
cp /app/go-fil-markets/config/lotusminer/config.toml ~/lotusminer/config.toml

# Start lotus node
LOTUS_DIR=/app/lotus ./run_lotus.sh &

# Start lotus miner
LOTUS_DIR=/app/lotus ./run_lotus_miner.sh &

# Start market
./run_market.sh