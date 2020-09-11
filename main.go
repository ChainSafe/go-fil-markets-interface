// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"github.com/ChainSafe/go-fil-markets-interface/retrievaladapter"
	"github.com/ChainSafe/go-fil-markets-interface/rpc"
	"github.com/ChainSafe/go-fil-markets-interface/storageadapter"
	"github.com/ChainSafe/go-fil-markets-interface/utils"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("markets")

func main() {
	lvl, err := logging.LevelFromString("INFO")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	flag.Parse()
	config.Load("./config/config.json")

	nodeClient, nodeCloser, err := nodeapi.GetNodeAPI(nil)
	if err != nil {
		log.Fatalf("Error while initializing Node client: %s", err)
	}
	log.Infof("Initialized node client")
	nodeapi.NodeClient = nodeClient

	params, err := utils.InitMarketParams()
	if err != nil {
		log.Fatalf("Error while initializing market client params: %s", err)
	}

	storageClient, err := storageadapter.InitStorageClient(params)
	if err != nil {
		log.Fatalf("Error while initializing storage client: %s", err)
	}
	log.Infof("Initialized storage market")

	retrievalClient, err := retrievaladapter.InitRetrievalClient(params)
	if err != nil {
		log.Fatalf("Error while initializing retrieval client: %s", err)
	}
	log.Infof("Initialized retrieval market")

	if err := rpc.Serve(storageClient, retrievalClient, params); err != nil {
		log.Fatalf("Error while setting up the server %s.", err)
	}
	log.Infof("Started serving Markets on %s", config.Api.Market.Addr)

	sdCh := make(chan os.Signal, 1)
	signal.Notify(sdCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan bool, 1)

	// handle signals
	go func() {
		var sigCnt int
		for sig := range sdCh {
			log.Infof("--- Received %s signal", sig)
			sigCnt++
			if sigCnt == 1 {
				// Graceful shutdown.
				signal.Stop(sdCh)
				doneCh <- true
				err := storageClient.Stop()
				if err != nil {
					log.Fatalf("Error while closing storage client %v", err)
				}
				_ = params.DataTransfer.Stop()
				nodeCloser()
			} else if sigCnt == 3 {
				// Force Shutdown
				log.Infof("--- Got interrupt signal 3rd time. Aborting now.")
				os.Exit(1)
			} else {
				log.Infof("--- Ignoring interrupt signal.")
			}
		}
	}()

	<-doneCh
}
