package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/crypto"
)

type Wallet struct{}

func (w *Wallet) Sign(ctx context.Context, addr address.Address, msg []byte) (*crypto.Signature, error) {
	return nil, nil
}

func (w *Wallet) GetDefault() (address.Address, error) {
	return address.Undef, nil
}
