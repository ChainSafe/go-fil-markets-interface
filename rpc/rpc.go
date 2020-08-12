// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package rpc

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChainSafe/go-fil-markets-interface/api"
	"github.com/ChainSafe/go-fil-markets-interface/auth"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/discovery"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-multistore"
	bstore "github.com/filecoin-project/lotus/lib/blockstore"
	"github.com/filecoin-project/lotus/node/repo/importmgr"
	"github.com/filecoin-project/lotus/node/repo/retrievalstoremgr"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	dss "github.com/ipfs/go-datastore/sync"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

var serverAddr = flag.String("server", "/ip4/127.0.0.1/tcp/7070/http", "server address")

func Serve(storageClient storagemarket.StorageClient, retrievalClient retrievalmarket.RetrievalClient) error {
	rpcServer := jsonrpc.NewServer()

	ds := dss.MutexWrap(datastore.NewMapDatastore())
	bs := bstore.NewBlockstore(namespace.Wrap(ds, datastore.NewKey("blockstore")))
	mds, err := multistore.NewMultiDstore(ds)
	if err != nil {
		return err
	}

	local := discovery.NewLocal(namespace.Wrap(ds, datastore.NewKey("/deals/local")))

	rpcServer.Register("Market", &api.API{
		RetDiscovery:      discovery.Multi(local),
		SMDealClient:      storageClient,
		Retrieval:         retrievalClient,
		CombinedBstore:    bs,
		Imports:           importmgr.New(mds, namespace.Wrap(ds, datastore.NewKey("/client"))),
		RetrievalStoreMgr: retrievalstoremgr.NewBlockstoreRetrievalStoreManager(bs),
	})

	ah := &auth.Handler{
		Next: rpcServer.ServeHTTP,
	}

	http.Handle("/rpc/market/v0/", ah)

	maddr, err := multiaddr.NewMultiaddr(*serverAddr)
	if err != nil {
		log.Fatalf("failed to construct multiaddr: %s %v", maddr, err)
	}

	lst, err := manet.Listen(maddr)
	if err != nil {
		return fmt.Errorf("could not listen: %w", err)
	}

	srv := &http.Server{Handler: http.DefaultServeMux}
	return srv.Serve(manet.NetListener(lst))
}
