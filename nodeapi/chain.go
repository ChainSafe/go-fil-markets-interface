// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
)

type ChainAPI struct {
	node Node
}

func (c *ChainAPI) ChainHead(ctx context.Context) (*types.TipSet, error) {
	return c.node.Chain.ChainHead(ctx)
}

func (c *ChainAPI) ChainGetTipSet(ctx context.Context, key types.TipSetKey) (*types.TipSet, error) {
	return c.node.Chain.ChainGetTipSet(ctx, key)
}

func (c *ChainAPI) ChainNotify(ctx context.Context) (<-chan []*api.HeadChange, error) {
	return c.node.Chain.ChainNotify(ctx)
}

func (c *ChainAPI) ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*api.BlockMessages, error) {
	return c.node.Chain.ChainGetBlockMessages(ctx, blockCid)
}

func (c *ChainAPI) ChainGetTipSetByHeight(ctx context.Context, e abi.ChainEpoch, tsk types.TipSetKey) (*types.TipSet, error) {
	return c.node.Chain.ChainGetTipSetByHeight(ctx, e, tsk)
}

type ApiBStore struct {
	node Node
}

func (a *ApiBStore) ChainReadObj(ctx context.Context, c cid.Cid) ([]byte, error) {
	return a.node.Chain.ChainReadObj(ctx, c)
}

func (a *ApiBStore) ChainHasObj(ctx context.Context, c cid.Cid) (bool, error) {
	return a.node.Chain.ChainHasObj(ctx, c)
}
