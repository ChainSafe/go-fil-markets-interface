// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import "context"

type MpoolAPI struct {}

// TODO(arijit): Implement the following to connect to Node and fetch info.
func (a *MpoolAPI) MpoolPushMessage(ctx context.Context) error {
	return nil
}

func (a *MpoolAPI) EnsureAvailable(ctx context.Context) error {
	return nil
}