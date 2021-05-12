package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

type FundManager struct{}

func (f *FundManager )Reserve(ctx context.Context, wallet, addr address.Address,
	amt abi.TokenAmount) (cid.Cid, error) {
	return NodeClient.FundManagerAPI.Reserve(ctx, wallet, addr, amt)
}

func (f *FundManager )Release(addr address.Address, amt abi.TokenAmount) error {
	return NodeClient.FundManagerAPI.Release(addr, amt)
}
