// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package rpc

import (
	"fmt"
	"net/http"

	"github.com/ChainSafe/go-fil-markets-interface/api"
	"github.com/ChainSafe/go-fil-markets-interface/auth"
	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/utils"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/multiformats/go-multiaddr/net"
)

func Serve(storageClient storagemarket.StorageClient, retrievalClient retrievalmarket.RetrievalClient, params *utils.MarketParams) error {
	rpcServer := jsonrpc.NewServer()
	marketAPI := &api.API{
		RetDiscovery:   params.Discovery,
		SMDealClient:   storageClient,
		Retrieval:      retrievalClient,
		CombinedBstore: params.Cbs,
		Imports:        modules.ClientImportMgr(params.Mds, params.Ds),
		Host:           params.Host,
		DataTransfer:   params.DataTransfer,
	}
	marketAPI.RetrievalStoreMgr = modules.ClientRetrievalStoreManager(marketAPI.Imports)

	rpcServer.Register("Market", marketAPI)

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
