// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package api

import (
	"context"
)

type RpcServer struct {}

func (s *RpcServer) GetString(ctx context.Context, token string) (string, error) {
	return token, nil
}