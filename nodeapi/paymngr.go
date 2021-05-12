// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"
)

type PaymentManager struct{}

func (pm *PaymentManager) GetPaych(ctx context.Context, from, to address.Address, ensureFree types.BigInt) (*api.ChannelInfo, error) {
	return NodeClient.PaymentManagerAPI.GetChannelInfo(ctx, from, to, ensureFree)
}

func (pm *PaymentManager) AllocateLane(ctx context.Context, ch address.Address) (uint64, error) {
	return NodeClient.PaymentManagerAPI.AllocateLane(ctx, ch)
}

// PaychVoucherCreate creates a new signed voucher on the given payment channel
// with the given lane and amount.  The value passed in is exactly the value
// that will be used to create the voucher, so if previous vouchers exist, the
// actual additional value of this voucher will only be the difference between
// the two.
func (pm *PaymentManager) PaychVoucherCreate(ctx context.Context, pch address.Address, amt types.BigInt, lane uint64) (*paych.SignedVoucher, error) {
	return NodeClient.PaymentManagerAPI.VoucherCreate(ctx, pch, amt, lane)
}

func (pm *PaymentManager) AvailableFunds(ctx context.Context, pch address.Address) (*api.ChannelAvailableFunds, error) {
	return NodeClient.PaymentManagerAPI.AvailableFunds(ctx, pch)
}

func (pm *PaymentManager) GetWaitReady(ctx context.Context, waitSentinel cid.Cid) (address.Address, error) {
	return NodeClient.PaymentManagerAPI.GetWaitReady(ctx, waitSentinel)
}