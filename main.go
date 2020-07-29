// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"flag"
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
	if err := rpc.Serve(); err != nil {
		log.Fatalf("Error while setting up the server.")
	}

	_ = storageadapter.NewStorageClientNode()
	_ = retrievaladapter.NewRetrievalClientNode()

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
