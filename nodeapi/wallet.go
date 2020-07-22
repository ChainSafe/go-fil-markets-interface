// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/crypto"
)

type Wallet struct {
	node *Client
}

func (w *Wallet) Sign(ctx context.Context, addr address.Address, msg []byte) (*crypto.Signature, error) {
	return w.node.Wallet.Sign(ctx, addr, msg)
}

func (w *Wallet) GetDefault() (address.Address, error) {
	return w.node.Wallet.GetDefault()
}
