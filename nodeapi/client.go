// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"
	"net/http"

	"github.com/ChainSafe/go-fil-markets-interface/config"

	"github.com/ChainSafe/go-fil-markets-interface/utils"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
)

type Node struct {
	Chain struct {
		ChainHead              func(ctx context.Context) (*types.TipSet, error)
		ChainGetTipSet         func(ctx context.Context, key types.TipSetKey) (*types.TipSet, error)
		ChainNotify            func(context.Context) (<-chan []*api.HeadChange, error)
		ChainGetBlockMessages  func(ctx context.Context, blockCid cid.Cid) (*api.BlockMessages, error)
		ChainGetTipSetByHeight func(ctx context.Context, e abi.ChainEpoch, ts types.TipSetKey) (*types.TipSet, error)
		ChainReadObj           func(context.Context, cid.Cid) ([]byte, error)
		ChainHasObj            func(context.Context, cid.Cid) (bool, error)
	}
	Mpool struct {
		MpoolPushMessage      func(ctx context.Context, msg *types.Message) (*types.SignedMessage, error)
		MarketEnsureAvailable func(ctx context.Context, addr, wallet address.Address, amt types.BigInt) (cid.Cid, error)
	}
	PaymentManager struct {
		PaychGet           func(ctx context.Context, from, to address.Address, ensureFree types.BigInt) (address.Address, cid.Cid, error)
		PaychAllocateLane  func(ctx context.Context, ch address.Address) (uint64, error)
		PaychVoucherCreate func(ctx context.Context, pch address.Address, amt types.BigInt, lane uint64) (*paych.SignedVoucher, error)
	}
	State struct {
		StateMarketBalance                func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error)
		StateAccountKey                   func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
		WaitForMessage                    func(ctx context.Context) error
		StateWaitMsg                      func(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error)
		StateMarketDeals                  func(ctx context.Context, tsk types.TipSetKey) (map[string]api.MarketDeal, error)
		StateListMiners                   func(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error)
		StateMinerInfo                    func(ctx context.Context, actor address.Address, tsk types.TipSetKey) (api.MinerInfo, error)
		StateLookupID                     func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
		StateMarketStorageDeal            func(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*api.MarketDeal, error)
		StateMinerProvingDeadline         func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*miner.DeadlineInfo, error)
		StateGetActor                     func(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error)
		StateGetReceipt                   func(context.Context, cid.Cid, types.TipSetKey) (*types.MessageReceipt, error)
		StateDealProviderCollateralBounds func(context.Context, abi.PaddedPieceSize, bool, types.TipSetKey) (api.DealCollateralBounds, error)
	}
	StateManager struct {
		StateWaitMsg  func(ctx context.Context, mcid cid.Cid, confidence uint64) (*api.MsgLookup, error)
		StateLookupID func(ctx context.Context, addr address.Address, ts types.TipSetKey) (address.Address, error)
	}
	ChainStore struct {
		ChainGetMessage func(ctx context.Context, c cid.Cid) (*types.Message, error)
		ChainHead       func(ctx context.Context) *types.TipSet
	}
	Wallet struct {
		WalletSign           func(ctx context.Context, addr address.Address, msg []byte) (*crypto.Signature, error)
		WalletDefaultAddress func(ctx context.Context) (address.Address, error)
		WalletHas            func(ctx context.Context, addr address.Address) (bool, error)
	}
	Utils struct {
		StateMarketStorageDeal func(ctx context.Context, dealID abi.DealID, ts types.TipSetKey) (*api.MarketDeal, error)
		StateMinerInfo         func(ctx context.Context, addr address.Address, ts types.TipSetKey) (api.MinerInfo, error)
	}
}

func NewNodeClient(addr string, requestHeader http.Header) (*Node, jsonrpc.ClientCloser, error) {
	var node Node
	closer, err := jsonrpc.NewMergeClient(addr, "MarketInterface",
		[]interface{}{
			&node.Chain,
			&node.Mpool,
			&node.PaymentManager,
			&node.State,
			&node.StateManager,
			&node.ChainStore,
			&node.Wallet,
			&node.Utils,
		},
		requestHeader)
	return &node, closer, err
}

func GetNodeAPI(ctx *cli.Context) (*Node, jsonrpc.ClientCloser, error) {
	addr, err := config.Api.Node.DialArgs()
	if err != nil {
		return nil, nil, err
	}
	return NewNodeClient(addr, utils.AuthHeader(string(config.Api.Node.Token)))
}
