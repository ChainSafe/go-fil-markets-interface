// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type Mpool struct{}

func (a *Mpool) MpoolPushMessage(ctx context.Context, msg *types.Message) (*types.SignedMessage, error) {
	return NodeClient.MpoolAPI.MpoolPushMessage(ctx, msg)
}

func (a *Mpool) EnsureAvailable(ctx context.Context, addr, wallet address.Address, amt types.BigInt) (cid.Cid, error) {
	return NodeClient.MpoolAPI.MarketEnsureAvailable(ctx, addr, wallet, amt)
}
