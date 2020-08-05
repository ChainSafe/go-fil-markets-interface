// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package retrievaladapter

import (
	"bytes"
	"context"

	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/lotus/build"
	initactor "github.com/filecoin-project/specs-actors/actors/builtin/init"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"golang.org/x/xerrors"

	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"
)

type ClientNodeAdapter struct {
	sm *nodeapi.StateManager

	nodeapi.ChainAPI
	nodeapi.PaymentManager
}

func NewRetrievalClientNode() retrievalmarket.RetrievalClientNode {
	return &ClientNodeAdapter{}
}

// GetChainHead gets the current chain head. Return its TipSetToken and its abi.ChainEpoch.
func (c *ClientNodeAdapter) GetChainHead(ctx context.Context) (shared.TipSetToken, abi.ChainEpoch, error) {
	head, err := c.ChainHead(ctx)
	if err != nil {
		return nil, 0, err
	}

	return head.Key().Bytes(), head.Height(), nil
}

// GetOrCreatePaymentChannel sets up a new payment channel if one does not exist
// between a client and a miner and ensures the client has the given amount of
// funds available in the channel.
func (c *ClientNodeAdapter) GetOrCreatePaymentChannel(ctx context.Context, clientAddress, minerAddress address.Address, clientFundsAvailable abi.TokenAmount, tok shared.TipSetToken) (address.Address, cid.Cid, error) {
	return c.GetPaych(ctx, clientAddress, minerAddress, clientFundsAvailable)
}

// Allocate late creates a lane within a payment channel so that calls to
// CreatePaymentVoucher will automatically make vouchers only for the difference
// in total
func (c *ClientNodeAdapter) AllocateLane(paymentChannel address.Address) (uint64, error) {
	return c.PaymentManager.AllocateLane(context.TODO(), paymentChannel)
}

// CreatePaymentVoucher creates a new payment voucher in the given lane for a
// given payment channel so that all the payment vouchers in the lane add up
// to the given amount (so the payment voucher will be for the difference)
func (c *ClientNodeAdapter) CreatePaymentVoucher(ctx context.Context, paymentChannel address.Address, amount abi.TokenAmount, lane uint64, tok shared.TipSetToken) (*paych.SignedVoucher, error) {
	voucher, err := c.PaychVoucherCreate(ctx, paymentChannel, amount, lane)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

// WaitForPaymentChannelAddFunds waits messageCID to appear on chain. If it doesn't appear within
// defaultMsgWaitTimeout it returns error
func (c *ClientNodeAdapter) WaitForPaymentChannelAddFunds(messageCID cid.Cid) error {
	_, mr, err := c.sm.WaitForMessage(context.TODO(), messageCID, build.MessageConfidence)

	if err != nil {
		return err
	}
	if mr.ExitCode != exitcode.Ok {
		return xerrors.Errorf("wait for payment channel to add funds failed. exit code: %d", mr.ExitCode)
	}
	return nil
}

// WaitForPaymentChannelCreation waits for a message on chain with CID messageCID that a payment channel has been created.
func (c *ClientNodeAdapter) WaitForPaymentChannelCreation(messageCID cid.Cid) (address.Address, error) {
	_, mr, err := c.sm.WaitForMessage(context.TODO(), messageCID, build.MessageConfidence)

	if err != nil {
		return address.Undef, err
	}
	if mr.ExitCode != exitcode.Ok {
		return address.Undef, xerrors.Errorf("payment channel creation failed. exit code: %d", mr.ExitCode)
	}
	var retval initactor.ExecReturn
	if err := retval.UnmarshalCBOR(bytes.NewReader(mr.Return)); err != nil {
		return address.Undef, err
	}
	return retval.RobustAddress, nil
}
