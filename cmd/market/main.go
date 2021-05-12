package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"github.com/ChainSafe/go-fil-markets-interface/retrievaladapter"
	"github.com/ChainSafe/go-fil-markets-interface/rpc"
	"github.com/ChainSafe/go-fil-markets-interface/storageadapter"
	"github.com/ChainSafe/go-fil-markets-interface/utils"
	lcli "github.com/filecoin-project/lotus/cli"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("market")

// DaemonCmd is the `go-fil-markets daemon` command
var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a Filecoin market daemon process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "./config/config.json",
		},
	},
	Action: func(cctx *cli.Context) error {
		config.Load(cctx.String("config"))

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
					_ = params.DataTransfer.Stop(context.Background())
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
		return nil
	},
}

func main() {
	lvl, err := logging.LevelFromString("INFO")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	local := []*cli.Command{
		DaemonCmd,
	}

	app := &cli.App{
		Name:                 "go-fil-market",
		Usage:                "Filecoin Market for storage and retrieval.",
		Version:              "0.1",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	lcli.RunApp(app)
}
