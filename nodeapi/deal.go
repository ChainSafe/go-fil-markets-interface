package nodeapi

import (
	"context"

	"github.com/filecoin-project/lotus/chain/actors/builtin/market"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
	"github.com/ipfs/go-cid"
)

type Deal struct{}

func (d *Deal) GetCurrentDealInfo(ctx context.Context, tok sealing.TipSetToken,
	proposal *market.DealProposal, publishCid cid.Cid) (sealing.CurrentDealInfo, error) {
	return NodeClient.DealAPI.GetCurrentDealInfo(ctx, tok, proposal, publishCid)
}
