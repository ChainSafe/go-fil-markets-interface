// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type StateManager struct{}

func (sm *StateManager) WaitForMessage(ctx context.Context, mcid cid.Cid, confidence uint64) (*types.TipSet, *types.MessageReceipt, error) {
	msg, err := NodeClient.StateManagerAPI.StateWaitMsg(ctx, mcid, confidence)
	if err != nil {
		return nil, nil, err
	}

	tipSet, err := NodeClient.ChainGetTipSet(context.TODO(), msg.TipSet)
	if err != nil {
		return nil, nil, err
	}
	return tipSet, &msg.Receipt, err
}

func (sm *StateManager) ResolveToKeyAddress(ctx context.Context, addr address.Address, ts *types.TipSet) (address.Address, error) {
	return NodeClient.StateManagerAPI.StateLookupID(ctx, addr, ts.Key())
}
