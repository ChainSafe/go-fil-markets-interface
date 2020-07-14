package nodeapi

import (
	"context"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

func GetStorageDeal(ctx context.Context, dealID abi.DealID, ts *types.TipSet) (*api.MarketDeal, error) {
	return &api.MarketDeal{}, nil
}
