// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package rpc

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChainSafe/fil-markets-interface/api"
	"github.com/ChainSafe/fil-markets-interface/auth"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

var serverAddr = flag.String("server", "/ip4/127.0.0.1/tcp/7070/http", "server address")

func Serve() error {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("MarketInterface", &api.RpcServer{})

	ah := &auth.Handler{
		Next: rpcServer.ServeHTTP,
	}

	http.Handle("/rpc", ah)

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
