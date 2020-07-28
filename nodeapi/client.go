package nodeapi

import (
	"context"
	"net/http"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/crypto"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
)

type Client struct {
	Chain struct {
		ChainHead      func(ctx context.Context) (*types.TipSet, error)
		ChainGetTipSet func(ctx context.Context, key types.TipSetKey) (*types.TipSet, error)
	}
	Mpool struct {
		PushMessage     func(ctx context.Context, msg *types.Message) (*types.SignedMessage, error)
		EnsureAvailable func(ctx context.Context, addr, wallet address.Address, amt types.BigInt) (cid.Cid, error)
	}
	PaymentManager struct {
		GetPaych           func(ctx context.Context, from, to address.Address, ensureFree types.BigInt) (address.Address, cid.Cid, error)
		AllocateLane       func(ch address.Address) (uint64, error)
		PaychVoucherCreate func(ctx context.Context, pch address.Address, amt types.BigInt, lane uint64) (*paych.SignedVoucher, error)
	}
	State struct {
		StateMarketBalance     func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error)
		StateAccountKey        func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
		WaitForMessage         func(ctx context.Context) error
		StateWaitMsg           func(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error)
		StateMarketDeals       func(ctx context.Context, tsk types.TipSetKey) (map[string]api.MarketDeal, error)
		StateListMiners        func(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error)
		StateMinerInfo         func(ctx context.Context, actor address.Address, tsk types.TipSetKey) (api.MinerInfo, error)
		StateLookupID          func(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
		StateMarketStorageDeal func(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*api.MarketDeal, error)
	}
	StateManager struct {
		WaitForMessage      func(ctx context.Context, mcid cid.Cid, confidence uint64) (*types.TipSet, *types.MessageReceipt, error)
		ResolveToKeyAddress func(ctx context.Context, addr address.Address, ts *types.TipSet) (address.Address, error)
	}
	ChainStore struct {
		GetMessage        func(c cid.Cid) (*types.Message, error)
		GetHeaviestTipSet func() *types.TipSet
	}
	Wallet struct {
		Sign       func(ctx context.Context, addr address.Address, msg []byte) (*crypto.Signature, error)
		GetDefault func() (address.Address, error)
	}
	Utils struct {
		GetStorageDeal func(ctx context.Context, dealID abi.DealID, ts *types.TipSet) (*api.MarketDeal, error)
		StateMinerInfo func(ctx context.Context, sm *StateManager, ts *types.TipSet, maddr address.Address) (miner.MinerInfo, error)
	}
}

func NewNodeClient(addr string, requestHeader http.Header) (*Client, jsonrpc.ClientCloser, error) {
	var node Client
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
