// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package utils

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	datatransfer "github.com/filecoin-project/go-data-transfer"
	dtimpl "github.com/filecoin-project/go-data-transfer/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/transport/graphsync"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	"github.com/filecoin-project/go-multistore"
	"github.com/filecoin-project/go-storedcounter"
	badgerbs "github.com/filecoin-project/lotus/blockstore/badger"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	badger "github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-graphsync"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/urfave/cli/v2"
)

var log = logging.Logger("utils")

type MarketParams struct {
	Host         host.Host
	Cbs          dtypes.ClientBlockstore
	Ds           dtypes.MetadataDS
	Mds          dtypes.ClientMultiDstore
	DataTransfer datatransfer.Manager
	Discovery    *discoveryimpl.Local
	Deals        dtypes.ClientDatastore
}

func AuthHeader(token string) http.Header {
	if len(token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(token))
		return headers
	}
	log.Warn("API Token not set and requested, capabilities might be limited.")
	return nil
}

func DaemonContext(cctx *cli.Context) context.Context {
	return context.Background()
}

func ReqContext(cctx *cli.Context) context.Context {
	tCtx := DaemonContext(cctx)

	ctx, done := context.WithCancel(tCtx)
	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		done()
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	return ctx
}

func NewLibP2PHost() (host.Host, error) {
	ctx := context.Background()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}
	return libp2p.New(ctx, opts...)
}

func InitMarketParams() (*MarketParams, error) {
	host, err := NewLibP2PHost()
	if err != nil {
		return nil, err
	}

	clientDir, err := ioutil.TempDir("", "client-ds")
	if err != nil {
		return nil, err
	}

	ds, err := badger.NewDatastore(clientDir, nil)
	if err != nil {
		return nil, err
	}

	mds, err := multistore.NewMultiDstore(ds)
	if err != nil {
		return nil, err
	}

	imgr := modules.ClientImportMgr(mds, namespace.Wrap(ds, datastore.NewKey("/client")))
	clientBs := modules.ClientBlockstore(imgr)

	loader := storeutil.LoaderForBlockstore(clientBs)
	storer := storeutil.StorerForBlockstore(clientBs)
	graphSyncNetwork := gsnet.NewFromLibp2pHost(host)

	chainDir, err := ioutil.TempDir("", "chain-ds")
	if err != nil {
		return nil, err
	}

	opts, err := repo.BadgerBlockstoreOptions(repo.HotBlockstore, chainDir, false)
	if err != nil {
		return nil, err
	}

	chainBs, err := badgerbs.Open(opts)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	graphSync := graphsyncimpl.New(ctx, graphSyncNetwork, loader, storer, graphsyncimpl.RejectAllRequestsByDefault())
	chainLoader := storeutil.LoaderForBlockstore(chainBs)
	chainStorer := storeutil.StorerForBlockstore(chainBs)

	err = graphSync.RegisterPersistenceOption("chainstore", chainLoader, chainStorer)
	if err != nil {
		return nil, err
	}

	graphSync.RegisterIncomingRequestHook(func(p peer.ID, requestData graphsync.RequestData, hookActions graphsync.IncomingRequestHookActions) {
		_, has := requestData.Extension("chainsync")
		if has {
			hookActions.ValidateRequest()
			hookActions.UsePersistenceOption("chainstore")
		}
	})
	graphSync.RegisterOutgoingRequestHook(func(p peer.ID, requestData graphsync.RequestData, hookActions graphsync.OutgoingRequestHookActions) {
		_, has := requestData.Extension("chainsync")
		if has {
			hookActions.UsePersistenceOption("chainstore")
		}
	})

	sc := storedcounter.New(ds, datastore.NewKey("/datatransfer/client/counter"))
	net := dtnet.NewFromLibp2pHost(host)

	dtDs := namespace.Wrap(ds, datastore.NewKey("/datatransfer/client/transfers"))
	transport := dtgstransport.NewTransport(host.ID(), graphSync)

	// data-transfer push channel restart configuration
	dtRestartConfig := dtimpl.PushChannelRestartConfig(time.Minute, 10, 1024, 10*time.Minute, 3)

	dtDir, err := ioutil.TempDir("", "data-transfer")
	if err != nil {
		return nil, err
	}
	dt, err := dtimpl.NewDataTransfer(dtDs, dtDir, net, transport, sc, dtRestartConfig)
	if err != nil {
		return nil, err
	}

	if err := dt.Start(ctx); err != nil {
		return nil, err
	}

	local, err := discoveryimpl.NewLocal(namespace.Wrap(ds, datastore.NewKey("/deals/local")))
	if err != nil {
		return nil, err
	}

	return &MarketParams{
		Host:         host,
		Cbs:          clientBs,
		Ds:           ds,
		Mds:          mds,
		DataTransfer: dt,
		Discovery:    local,
		Deals:        modules.NewClientDatastore(ds),
	}, nil
}
