// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package retrievaladapter

import (
	"context"

	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	mutils "github.com/ChainSafe/go-fil-markets-interface/utils"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/discovery"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	retrievalimpl "github.com/filecoin-project/go-fil-markets/retrievalmarket/impl"
	rmnet "github.com/filecoin-project/go-fil-markets/retrievalmarket/network"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-storedcounter"
	"github.com/filecoin-project/lotus/chain/types"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/markets/retrievaladapter"
	"github.com/filecoin-project/lotus/node/impl/full"
	payapi "github.com/filecoin-project/lotus/node/impl/paych"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)

type ClientNodeAdapter struct {
	nodeapi.StateManager
	nodeapi.Chain
	nodeapi.PaymentManager
	nodeapi.State
}

func InitRetrievalClient(params *mutils.MarketParams) (retrievalmarket.RetrievalClient, error) {
	rcn := NewRetrievalClientNode()
	retrievalClient, err := RetrievalClient(params.Host, params.Mds, params.DataTransfer, params.Discovery, params.Ds, rcn)
	if err != nil {
		return nil, err
	}
	return retrievalClient, err
}

func RetrievalClient(h host.Host, mds dtypes.ClientMultiDstore, dt dtypes.ClientDataTransfer, resolver discovery.PeerResolver, ds dtypes.MetadataDS, adapter retrievalmarket.RetrievalClientNode) (retrievalmarket.RetrievalClient, error) {
	network := rmnet.NewFromLibp2pHost(h)
	sc := storedcounter.New(ds, datastore.NewKey("/retr"))
	client, err := retrievalimpl.NewClient(network, mds, dt, adapter, resolver, namespace.Wrap(ds, datastore.NewKey("/retrievals/client")), sc)
	if err != nil {
		return nil, err
	}
	client.SubscribeToEvents(marketevents.RetrievalClientLogger)
	return client, nil
}

func NewRetrievalClientNode() retrievalmarket.RetrievalClientNode {
	return retrievaladapter.NewRetrievalClientNode(payapi.PaychAPI{}, full.ChainAPI{}, full.StateAPI{})
}

func (c *ClientNodeAdapter) WaitForPaymentChannelReady(ctx context.Context, waitSentinel cid.Cid) (address.Address, error) {
	return c.PaymentManager.GetWaitReady(ctx, waitSentinel)
}

func (c *ClientNodeAdapter) CheckAvailableFunds(ctx context.Context, paymentChannel address.Address) (retrievalmarket.ChannelAvailableFunds, error) {
	channelAvailableFunds, err := c.AvailableFunds(ctx, paymentChannel)
	if err != nil {
		return retrievalmarket.ChannelAvailableFunds{}, err
	}
	return retrievalmarket.ChannelAvailableFunds{
		ConfirmedAmt:        channelAvailableFunds.ConfirmedAmt,
		PendingAmt:          channelAvailableFunds.PendingAmt,
		PendingWaitSentinel: channelAvailableFunds.PendingWaitSentinel,
		QueuedAmt:           channelAvailableFunds.QueuedAmt,
		VoucherReedeemedAmt: channelAvailableFunds.VoucherReedeemedAmt,
	}, nil
}

func (c *ClientNodeAdapter) GetKnownAddresses(ctx context.Context, p retrievalmarket.RetrievalPeer, tok shared.TipSetToken) ([]multiaddr.Multiaddr, error) {
	tsk, err := types.TipSetKeyFromBytes(tok)
	if err != nil {
		return nil, err
	}
	mi, err := c.StateMinerInfo(ctx, p.Address, tsk)
	if err != nil {
		return nil, err
	}
	multiaddrs := make([]multiaddr.Multiaddr, 0, len(mi.Multiaddrs))
	for _, a := range mi.Multiaddrs {
		maddr, err := multiaddr.NewMultiaddrBytes(a)
		if err != nil {
			return nil, err
		}
		multiaddrs = append(multiaddrs, maddr)
	}

	return multiaddrs, nil
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
	chanInfo, err := c.GetPaych(ctx, clientAddress, minerAddress, clientFundsAvailable)
	return chanInfo.Channel, chanInfo.WaitSentinel, err
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