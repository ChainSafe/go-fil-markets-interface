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
	node   *Node
}

func (s *StateAPI) StateMarketBalance(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error) {
	return s.node.State.StateMarketBalance(ctx, addr, tsk)
}

func (s *StateAPI) StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return s.node.State.StateAccountKey(ctx, addr, tsk)
}

func (s *StateAPI) StateWaitMsg(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error) {
	return s.node.State.StateWaitMsg(ctx, msg, confidence)
}

func (s *StateAPI) StateMarketDeals(ctx context.Context, tsk types.TipSetKey) (map[string]api.MarketDeal, error) {
	return s.node.State.StateMarketDeals(ctx, tsk)
}

func (s *StateAPI) StateListMiners(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error) {
	return s.node.State.StateListMiners(ctx, tsk)
}

func (s *StateAPI) StateMinerInfo(ctx context.Context, actor address.Address, tsk types.TipSetKey) (api.MinerInfo, error) {
	return s.node.State.StateMinerInfo(ctx, actor, tsk)
}

func (s *StateAPI) StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return s.node.State.StateLookupID(ctx, addr, tsk)
}

func (s *StateAPI) StateMarketStorageDeal(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*api.MarketDeal, error) {
	return s.node.State.StateMarketStorageDeal(ctx, dealId, tsk)
}

func (s *StateAPI) StateMinerProvingDeadline(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*miner.DeadlineInfo, error) {
	return s.node.State.StateMinerProvingDeadline(ctx, addr, tsk)
}

func (s *StateAPI) StateGetReceipt(ctx context.Context, cid cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error) {
	return s.node.State.StateGetReceipt(ctx, cid, tsk)
}

func (s *StateAPI) StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) {
	return s.node.State.StateGetActor(ctx, actor, tsk)
}

func (s *StateAPI) StateDealProviderCollateralBounds(ctx context.Context, size abi.PaddedPieceSize, verified bool, tsk types.TipSetKey) (api.DealCollateralBounds, error) {
	return s.node.State.StateDealProviderCollateralBounds(ctx, size, verified, tsk)
}
