// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package api

import (
	"context"
)

type RpcServer struct{}

// Demo RPC API.
func (s *RpcServer) GetString(ctx context.Context, token string) (string, error) {
	return token, nil
}
