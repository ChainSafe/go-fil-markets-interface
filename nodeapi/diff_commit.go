package nodeapi

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
)

type DiffPreCommitsAPI struct{}

func (d *DiffPreCommitsAPI) DiffPreCommits(ctx context.Context, actor address.Address,
	pre, cur types.TipSetKey) (*miner.PreCommitChanges, error) {
	return NodeClient.CommitsAPI.DiffPreCommits(ctx, actor, pre, cur)
}
