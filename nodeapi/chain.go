// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
)

type ChainAPI struct {
	node *Node
}

func (a *ChainAPI) ChainHead(ctx context.Context) (*types.TipSet, error) {
	return a.node.Chain.ChainHead(ctx)
}

func (a *ChainAPI) ChainGetTipSet(ctx context.Context, key types.TipSetKey) (*types.TipSet, error) {
	return a.node.Chain.ChainGetTipSet(ctx, key)
}
