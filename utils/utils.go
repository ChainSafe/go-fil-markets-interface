package utils

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	datatransfer "github.com/filecoin-project/go-data-transfer"
	dtimpl "github.com/filecoin-project/go-data-transfer/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/transport/graphsync"
	"github.com/filecoin-project/go-data-transfer/transport/graphsync/extension"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/discovery"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/funds"
	"github.com/filecoin-project/go-multistore"
	"github.com/filecoin-project/go-storedcounter"
	bstore "github.com/filecoin-project/lotus/lib/blockstore"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo/importmgr"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	badger "github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-graphsync"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
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
	Ds           datastore.Batching
	Mds          *multistore.MultiStore
	DataTransfer datatransfer.Manager
	Discovery    *discovery.Local
	Deals        dtypes.ClientDatastore
	DealFunds    funds.DealFunds
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

func NewClientBlockStore(mds dtypes.ClientMultiDstore, ds dtypes.MetadataDS) dtypes.ChainBlockstore {
	imgr := importmgr.New(mds, namespace.Wrap(ds, datastore.NewKey("/client")))
	return bstore.WrapIDStore(imgr.Blockstore)
}

func NewChainBlockStore(ds dtypes.MetadataDS) (dtypes.ChainBlockstore, error) {
	bs := blockstore.NewBlockstore(ds)
	cbs, err := blockstore.CachedBlockstore(context.Background(), bs, blockstore.DefaultCacheOpts())
	if err != nil {
		return nil, err
	}

	return cbs, nil
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

	ctx := context.Background()
	tdir, err := ioutil.TempDir("", "market-client")
	if err != nil {
		return nil, err
	}

	ds, err := badger.NewDatastore(tdir, nil)
	if err != nil {
		return nil, err
	}

	mds, err := multistore.NewMultiDstore(ds)
	if err != nil {
		return nil, err
	}

	clientBs := NewClientBlockStore(mds, ds)

	loader := storeutil.LoaderForBlockstore(clientBs)
	storer := storeutil.StorerForBlockstore(clientBs)
	graphSyncNetwork := gsnet.NewFromLibp2pHost(host)

	chainBs, err := NewChainBlockStore(ds)
	if err != nil {
		return nil, err
	}

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
			// TODO: we should confirm the selector is a reasonable one before we validate
			// TODO: this code will get more complicated and should probably not live here eventually
			hookActions.ValidateRequest()
			hookActions.UsePersistenceOption("chainstore")
		}
		_, has = requestData.Extension(extension.ExtensionDataTransfer)
		if has {
			hookActions.ValidateRequest()
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
	dt, err := dtimpl.NewDataTransfer(dtDs, net, transport, sc)
	if err != nil {
		return nil, err
	}

	if err := dt.Start(ctx); err != nil {
		return nil, err
	}

	clientDealFunds, err := modules.NewClientDealFunds(ds)
	if err != nil {
		return nil, err
	}

	return &MarketParams{
		Host:         host,
		Cbs:          clientBs,
		Ds:           ds,
		Mds:          mds,
		DataTransfer: dt,
		Discovery:    modules.NewLocalDiscovery(ds),
		Deals:        modules.NewClientDatastore(ds),
		DealFunds:    clientDealFunds,
	}, nil
}
