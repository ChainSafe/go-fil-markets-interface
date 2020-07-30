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
	"github.com/ipfs/go-cid"
)

type StateAPI struct {
	Wallet Wallet
	node   *Client
}

func (a *StateAPI) StateMarketBalance(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error) {
	return a.node.State.StateMarketBalance(ctx, addr, tsk)
}

func (a *StateAPI) StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return a.node.State.StateAccountKey(ctx, addr, tsk)
}

func (a *StateAPI) WaitForMessage(ctx context.Context) error {
	return a.node.State.WaitForMessage(ctx)
}

// Keep polling till the Msg is received.
func (a *StateAPI) StateWaitMsg(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error) {
	return a.node.State.StateWaitMsg(ctx, msg, confidence)
}

func (a *StateAPI) StateMarketDeals(ctx context.Context, tsk types.TipSetKey) (map[string]api.MarketDeal, error) {
	return a.node.State.StateMarketDeals(ctx, tsk)
}

func (a *StateAPI) StateListMiners(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error) {
	return a.node.State.StateListMiners(ctx, tsk)
}

func (a *StateAPI) StateMinerInfo(ctx context.Context, actor address.Address, tsk types.TipSetKey) (api.MinerInfo, error) {
	return a.node.State.StateMinerInfo(ctx, actor, tsk)
}

func (a *StateAPI) StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return a.node.State.StateLookupID(ctx, addr, tsk)
}

func (a *StateAPI) StateMarketStorageDeal(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*api.MarketDeal, error) {
	return a.node.State.StateMarketStorageDeal(ctx, dealId, tsk)
}

func (a *StateAPI) StateMinerProvingDeadline(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*miner.DeadlineInfo, error) {
	return a.node.State.StateMinerProvingDeadline(ctx, addr, tsk)
}
