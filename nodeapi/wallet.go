// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api"
)

type Wallet struct{}

func (w *Wallet) GetDefault() (address.Address, error) {
	return NodeClient.WalletAPI.WalletDefaultAddress(context.TODO())
}

func (w *Wallet) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	return NodeClient.WalletAPI.WalletHas(ctx, addr)
}

func (w *Wallet) Sign(ctx context.Context, signer address.Address, toSign []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	return NodeClient.WalletAPI.WalletSign(ctx, signer, toSign, meta)
}
