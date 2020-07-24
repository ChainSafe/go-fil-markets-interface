// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type ChainStore struct {
	node *Client
}

func (cs *ChainStore) GetMessage(c cid.Cid) (*types.Message, error) {
	return cs.node.ChainStore.GetMessage(c)
}

func (cs *ChainStore) GetHeaviestTipSet() *types.TipSet {
	return cs.node.ChainStore.GetHeaviestTipSet()
}
