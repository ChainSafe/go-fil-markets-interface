// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

func GetStorageDeal(ctx context.Context, client Node, dealID abi.DealID, ts *types.TipSet) (*api.MarketDeal, error) {
	return NodeClient.UtilsAPI.StateMarketStorageDeal(ctx, dealID, ts.Key())
}

func StateMinerInfo(ctx context.Context, client Node, ts *types.TipSet, maddr address.Address) (api.MinerInfo, error) {
	return NodeClient.UtilsAPI.StateMinerInfo(ctx, maddr, ts.Key())
}
