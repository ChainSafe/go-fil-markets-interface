// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type ChainStore struct{}

func (cs *ChainStore) GetMessage(c cid.Cid) (*types.Message, error) {
	return NodeClient.ChainStoreAPI.ChainGetMessage(context.TODO(), c)
}

func (cs *ChainStore) GetHeaviestTipSet() *types.TipSet {
	return NodeClient.ChainStoreAPI.ChainHead(context.TODO())
}
