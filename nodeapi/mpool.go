// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type MpoolAPI struct {
	node *Node
}

func (a *MpoolAPI) MpoolPushMessage(ctx context.Context, msg *types.Message) (*types.SignedMessage, error) {
	return a.node.Mpool.PushMessage(ctx, msg)
}

func (a *MpoolAPI) EnsureAvailable(ctx context.Context, addr, wallet address.Address, amt types.BigInt) (cid.Cid, error) {
	return a.node.Mpool.EnsureAvailable(ctx, addr, wallet, amt)
}
