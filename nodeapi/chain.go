// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type Chain struct{}

func (c Chain) ChainHead(ctx context.Context) (*types.TipSet, error) {
	return NodeClient.ChainAPI.ChainHead(ctx)
}

func (c Chain) ChainGetTipSet(ctx context.Context, key types.TipSetKey) (*types.TipSet, error) {
	return NodeClient.ChainAPI.ChainGetTipSet(ctx, key)
}

func (c Chain) ChainNotify(ctx context.Context) (<-chan []*api.HeadChange, error) {
	return NodeClient.ChainAPI.ChainNotify(ctx)
}

func (c Chain) ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*api.BlockMessages, error) {
	return NodeClient.ChainAPI.ChainGetBlockMessages(ctx, blockCid)
}

func (c Chain) ChainGetTipSetByHeight(ctx context.Context, e abi.ChainEpoch, tsk types.TipSetKey) (*types.TipSet, error) {
	return NodeClient.ChainAPI.ChainGetTipSetByHeight(ctx, e, tsk)
}

type ApiBStore struct{}

func (a *ApiBStore) ChainReadObj(ctx context.Context, c cid.Cid) ([]byte, error) {
	return NodeClient.ChainAPI.ChainReadObj(ctx, c)
}

func (a *ApiBStore) ChainHasObj(ctx context.Context, c cid.Cid) (bool, error) {
	return NodeClient.ChainAPI.ChainHasObj(ctx, c)
}
