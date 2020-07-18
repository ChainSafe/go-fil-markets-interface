// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
)

func GetStorageDeal(ctx context.Context, dealID abi.DealID, ts *types.TipSet) (*api.MarketDeal, error) {
	return &api.MarketDeal{}, nil
}

func StateMinerInfo(ctx context.Context, sm *StateManager, ts *types.TipSet, maddr address.Address) (miner.MinerInfo, error) {
	return miner.MinerInfo{}, nil
}
