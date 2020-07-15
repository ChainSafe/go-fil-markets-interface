// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
)

type ChainAPI struct{}

// TODO(arijit): Implement the following to connect to Node and fetch info.
func (a *ChainAPI) ChainHead(context.Context) (*types.TipSet, error) {
	return nil, nil
}

func (a *ChainAPI) ChainGetTipSet(ctx context.Context, key types.TipSetKey) (*types.TipSet, error) {
	return nil, nil
}
