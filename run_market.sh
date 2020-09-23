#!/bin/sh

set -o xtrace

source ./init_market_env.sh
go build -gcflags='-N -l' -tags=2k
./go-fil-markets-interface
