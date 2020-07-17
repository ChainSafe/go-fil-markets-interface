// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"
)

type PaymentManager struct{}

func (pm *PaymentManager) GetPaych(ctx context.Context, from, to address.Address, ensureFree types.BigInt) (address.Address, cid.Cid, error) {
	return address.Undef, cid.Undef, nil
}

func (pm *PaymentManager) AllocateLane(ch address.Address) (uint64, error) {
	return 0, nil
}

// PaychVoucherCreate creates a new signed voucher on the given payment channel
// with the given lane and amount.  The value passed in is exactly the value
// that will be used to create the voucher, so if previous vouchers exist, the
// actual additional value of this voucher will only be the difference between
// the two.
func (pm *PaymentManager) PaychVoucherCreate(ctx context.Context, pch address.Address, amt types.BigInt, lane uint64) (*paych.SignedVoucher, error) {
	return nil, nil
}