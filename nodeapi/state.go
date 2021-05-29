// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/dline"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type State struct{}

func (s State) StateMarketBalance(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error) {
	return NodeClient.StateAPI.StateMarketBalance(ctx, addr, tsk)
}

func (s State) StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return NodeClient.StateAPI.StateAccountKey(ctx, addr, tsk)
}

func (s State) StateWaitMsg(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error) {
	return NodeClient.StateAPI.StateWaitMsg(ctx, msg, confidence)
}

func (s State) StateMarketDeals(ctx context.Context, tsk types.TipSetKey) (map[string]api.MarketDeal, error) {
	return NodeClient.StateAPI.StateMarketDeals(ctx, tsk)
}

func (s State) StateListMiners(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error) {
	return NodeClient.StateAPI.StateListMiners(ctx, tsk)
}

func (s State) StateMinerInfo(ctx context.Context, actor address.Address, tsk types.TipSetKey) (miner.MinerInfo, error) {
	return NodeClient.StateAPI.StateMinerInfo(ctx, actor, tsk)
}

func (s State) StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return NodeClient.StateAPI.StateLookupID(ctx, addr, tsk)
}

func (s State) StateMarketStorageDeal(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*api.MarketDeal, error) {
	return NodeClient.StateAPI.StateMarketStorageDeal(ctx, dealId, tsk)
}

func (s State) StateMinerProvingDeadline(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*dline.Info, error) {
	return NodeClient.StateAPI.StateMinerProvingDeadline(ctx, addr, tsk)
}

func (s State) StateGetReceipt(ctx context.Context, cid cid.Cid, tsk types.TipSetKey) (*types.MessageReceipt, error) {
	return NodeClient.StateAPI.StateGetReceipt(ctx, cid, tsk)
}

func (s State) StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) {
	return NodeClient.StateAPI.StateGetActor(ctx, actor, tsk)
}

func (s State) StateDealProviderCollateralBounds(ctx context.Context, size abi.PaddedPieceSize, verified bool, tsk types.TipSetKey) (api.DealCollateralBounds, error) {
	return NodeClient.StateAPI.StateDealProviderCollateralBounds(ctx, size, verified, tsk)
}

func (s State) StateNetworkVersion(ctx context.Context, key types.TipSetKey) (network.Version, error) {
	return NodeClient.StateAPI.StateNetworkVersion(ctx, key)
}