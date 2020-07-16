package nodeapi

import (
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type ChainStore struct{}

func (cs *ChainStore) GetMessage(c cid.Cid) (*types.Message, error) {
	return nil, nil
}

func (cs *ChainStore) GetHeaviestTipSet() *types.TipSet {
	return nil
}
