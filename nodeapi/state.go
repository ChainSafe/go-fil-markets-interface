// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type StateAPI struct{}

// TODO(arijit): Implement the following to connect to Node and fetch info.
func (a *StateAPI) StateMarketBalance(ctx context.Context, addr address.Address, tsk types.TipSetKey) (api.MarketBalance, error) {
	return api.MarketBalance{}, nil
}

func (a *StateAPI) StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error) {
	return address.Undef, nil
}

func (a *StateAPI) WaitForMessage(ctx context.Context) error {
	return nil
}

// Keep polling till the Msg is received.
func (a *StateAPI) StateWaitMsg(ctx context.Context, msg cid.Cid, confidence uint64) (*api.MsgLookup, error) {
	return nil, nil
}
