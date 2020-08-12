// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"flag"
	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/go-fil-markets-interface/retrievaladapter"
	"github.com/ChainSafe/go-fil-markets-interface/rpc"
	"github.com/ChainSafe/go-fil-markets-interface/storageadapter"
)

func main() {
	flag.Parse()
	config.Load("./config/config.json")

	nodeClient, nodeCloser, err := nodeapi.GetNodeAPI(nil)
	if err != nil {
		log.Fatalf("Error while initializing Node client: %s", err)
	}

	nodeapi.NodeClient = nodeClient

	storageClient, err := storageadapter.InitStorageClient(nodeClient)
	if err != nil {
		log.Fatalf("Error while initializing storage client: %s", err)
	}
	retrievalClient, err := retrievaladapter.InitRetrievalClient()
	if err != nil {
		log.Fatalf("Error while initializing retrieval client: %s", err)
	}

	if err := rpc.Serve(storageClient, retrievalClient); err != nil {
		log.Fatalf("Error while setting up the server.")
	}

	sdCh := make(chan os.Signal, 1)
	signal.Notify(sdCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan bool, 1)

	// handle signals
	go func() {
		var sigCnt int
		for sig := range sdCh {
			log.Printf("--- Received %s signal", sig)
			sigCnt++
			if sigCnt == 1 {
				// Graceful shutdown.
				signal.Stop(sdCh)
				doneCh <- true
				err := storageClient.Stop()
				if err != nil {
					log.Fatalf("Error while closing storage client %v", err)
				}
				nodeCloser()
			} else if sigCnt == 3 {
				// Force Shutdown
				log.Printf("--- Got interrupt signal 3rd time. Aborting now.")
				os.Exit(1)
			} else {
				log.Printf("--- Ignoring interrupt signal.")
			}
		}
	}()

	<-doneCh
}
