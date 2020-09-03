package utils

import (
	"context"
	bstore "github.com/filecoin-project/lotus/lib/blockstore"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo/importmgr"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

var log = logging.Logger("utils")

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
