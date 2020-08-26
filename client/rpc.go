package client

import (
	"context"
	"io"
	"net/http"

	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/utils"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/urfave/cli/v2"
)

type RepoType int

type Market struct {
	ClientStartDeal       func(ctx context.Context, params *api.StartDealParams) (*cid.Cid, error)
	ClientListDeals       func(ctx context.Context) ([]api.DealInfo, error)
	ClientGetDealInfo     func(ctx context.Context, d cid.Cid) (*api.DealInfo, error)
	ClientHasLocal        func(ctx context.Context, root cid.Cid) (bool, error)
	ClientFindData        func(ctx context.Context, root cid.Cid, piece *cid.Cid) ([]api.QueryOffer, error)
	ClientMinerQueryOffer func(ctx context.Context, miner address.Address, root cid.Cid, piece *cid.Cid) (api.QueryOffer, error)
	ClientImport          func(ctx context.Context, ref api.FileRef) (*api.ImportRes, error)
	ClientRemoveImport    func(ctx context.Context, importID int) error
	ClientImportLocal     func(ctx context.Context, f io.Reader) (cid.Cid, error)
	ClientListImports     func(ctx context.Context) ([]api.Import, error)
	ClientRetrieve        func(ctx context.Context, order api.RetrievalOrder, ref *api.FileRef) error
	ClientQueryAsk        func(ctx context.Context, p peer.ID, miner address.Address) (*storagemarket.SignedStorageAsk, error)
	ClientCalcCommP       func(ctx context.Context, inpath string, miner address.Address) (*api.CommPRet, error)
	ClientGenCar          func(ctx context.Context, ref api.FileRef, outputPath string) error
}

func NewMarketClient(addr string, requestHeader http.Header) (*Market, jsonrpc.ClientCloser, error) {
	var market Market
	closer, err := jsonrpc.NewMergeClient(context.Background(), addr, "Market",
		[]interface{}{
			&market,
		},
		requestHeader)
	return &market, closer, err
}

func GetMarketAPI(ctx *cli.Context) (*Market, jsonrpc.ClientCloser, error) {
	addr, err := config.Api.Market.DialArgs()
	if err != nil {
		return nil, nil, err
	}
	return NewMarketClient(addr, utils.AuthHeader(string(config.Api.Market.Token)))
}
