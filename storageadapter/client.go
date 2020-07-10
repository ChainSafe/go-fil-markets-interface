// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package storageadapter

import (
	"context"
	"github.com/ChainSafe/fil-markets-interface/nodeapi"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/ipfs/go-cid"
)

// This file implements StorageClientNode which is a client interface for making storage deals
// with a StorageProvider

type ClientNodeAdapter struct {
	nodeapi.ChainAPI
	nodeapi.MpoolAPI
	nodeapi.StateAPI
}

func NewStorageCommonImpl() storagemarket.StorageCommon {
	return &ClientNodeAdapter{}
}

// GetChainHead returns a tipset token for the current chain head
func (n *ClientNodeAdapter) GetChainHead(ctx context.Context) (shared.TipSetToken, abi.ChainEpoch, error) {
	err := n.ChainHead(ctx)
	if err != nil {
		return nil, 0, err
	}

	return nil, 0, nil
}

// Adds funds with the StorageMinerActor for a storage participant.  Used by both providers and clients.
func (n *ClientNodeAdapter) AddFunds(ctx context.Context, addr address.Address, amount abi.TokenAmount) (cid.Cid, error) {
	// (Provider Node API)
	err := n.MpoolPushMessage(ctx)
	return cid.Cid{}, err
}


// EnsureFunds ensures that a storage market participant has a certain amount of available funds
// If additional funds are needed, they will be sent from the 'wallet' address, and a cid for the
// corresponding chain message is returned
func (n *ClientNodeAdapter) EnsureFunds(ctx context.Context, addr, wallet address.Address, amount abi.TokenAmount, tok shared.TipSetToken) (cid.Cid, error) {
	err := n.MpoolPushMessage(ctx)
	return cid.Cid{}, err
}

// GetBalance returns locked/unlocked for a storage participant.  Used by both providers and clients.
func (n *ClientNodeAdapter) GetBalance(ctx context.Context, addr address.Address, tok shared.TipSetToken) (storagemarket.Balance, error) {
	err := n.StateMarketBalance(ctx)
	return storagemarket.Balance{}, err
}

// VerifySignature verifies a given set of data was signed properly by a given address's private key
func (n *ClientNodeAdapter) VerifySignature(ctx context.Context, signature crypto.Signature, signer address.Address, plaintext []byte, tok shared.TipSetToken) (bool, error) {
	err := n.StateAccountKey(ctx)
	return err != nil, err
}

// WaitForMessage waits until a message appears on chain. If it is already on chain, the callback is called immediately
func (n *ClientNodeAdapter) WaitForMessage(ctx context.Context, mcid cid.Cid, onCompletion func(exitcode.ExitCode, []byte, error) error) error {
	err := n.StateAccountKey(ctx)
	return err
}

// SignsBytes signs the given data with the given address's private key
func (n *ClientNodeAdapter) SignBytes(ctx context.Context, signer address.Address, b []byte) (*crypto.Signature, error) {
	return nil, nil
}

// OnDealSectorCommitted waits for a deal's sector to be sealed and proved, indicating the deal is active
func (n *ClientNodeAdapter) OnDealSectorCommitted(ctx context.Context, provider address.Address, dealID abi.DealID, cb storagemarket.DealSectorCommittedCallback) error {
	return nil
}