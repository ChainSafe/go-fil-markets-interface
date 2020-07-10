// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"
)

type StateAPI struct {}

// TODO(arijit): Implement the following to connect to Node and fetch info.
func (a *StateAPI) StateMarketBalance(ctx context.Context) error {
	return nil
}

func (a *StateAPI) StateAccountKey(ctx context.Context) error {
	return nil
}

func (a *StateAPI) WaitForMessage(ctx context.Context) error {
	return nil
}