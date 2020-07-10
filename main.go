package main

import (
	"flag"
	"github.com/ChainSafe/fil-markets-interface/rpc"
	"github.com/ChainSafe/fil-markets-interface/storageadapter"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	if err := rpc.Serve(); err != nil {
		log.Fatalf("Error while setting up the server.")
	}

	_ = storageadapter.ClientNodeAdapter{}

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