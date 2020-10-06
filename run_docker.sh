#!/usr/bin/env bash

make submodule

# Start lotus node
LOTUS_DIR=/app/lotus ./run_lotus.sh &

# Start lotus miner
LOTUS_DIR=/app/lotus ./run_lotus_miner.sh &

# Start market
./run_market.sh
