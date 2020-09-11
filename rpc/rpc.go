// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package rpc

import (
	"fmt"
	"net/http"

	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/utils"

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
	"github.com/multiformats/go-multiaddr/net"
)

func Serve(storageClient storagemarket.StorageClient, retrievalClient retrievalmarket.RetrievalClient, params *utils.MarketClientParams) error {
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
		Host:              params.Host,
		DataTransfer:      params.DataTransfer,
	})

	ah := &auth.Handler{
		Next: rpcServer.ServeHTTP,
	}

	http.Handle("/rpc/v0", ah)

	lst, err := manet.Listen(config.Api.Market.Addr)
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}

	srv := &http.Server{Handler: http.DefaultServeMux}
	return srv.Serve(manet.NetListener(lst))
}
