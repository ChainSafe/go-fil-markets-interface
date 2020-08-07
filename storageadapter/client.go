// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package storageadapter

import (
	"bytes"
	"context"
	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/events"
	"github.com/filecoin-project/lotus/chain/events/state"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/sigs"
	"github.com/filecoin-project/lotus/markets/utils"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	samarket "github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/golang/glog"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// This file implements StorageClientNode which is a client interface for making storage deals
// with a StorageProvider

type ClientNodeAdapter struct {
	cs *nodeapi.ChainStore
	sm *nodeapi.StateManager

	nodeapi.ChainAPI
	nodeapi.MpoolAPI
	nodeapi.StateAPI
	nodeapi.ApiBStore

	node *nodeapi.Node
	ev   *events.Events
}

type clientApi struct {
	nodeapi.ChainAPI
	nodeapi.StateAPI
}

func NewStorageClientNode() storagemarket.StorageClientNode {
	return &ClientNodeAdapter{
		ev: events.NewEvents(context.TODO(), &clientApi{nodeapi.ChainAPI{}, nodeapi.StateAPI{}}),
	}
}

func (n *ClientNodeAdapter) DealProviderCollateralBounds(ctx context.Context, size abi.PaddedPieceSize, isVerified bool) (abi.TokenAmount, abi.TokenAmount, error) {
	bounds, err := n.StateDealProviderCollateralBounds(ctx, size, isVerified, types.EmptyTSK)
	if err != nil {
		return abi.TokenAmount{}, abi.TokenAmount{}, err
	}

	return bounds.Min, bounds.Max, nil
}

func (n *ClientNodeAdapter) ListClientDeals(ctx context.Context, addr address.Address, encodedTs shared.TipSetToken) ([]storagemarket.StorageDeal, error) {
	tsk, err := types.TipSetKeyFromBytes(encodedTs)
	if err != nil {
		return nil, err
	}

	allDeals, err := n.StateMarketDeals(ctx, tsk)
	if err != nil {
		return nil, err
	}

	var out []storagemarket.StorageDeal

	for _, deal := range allDeals {
		storageDeal := utils.FromOnChainDeal(deal.Proposal, deal.State)
		if storageDeal.Client == addr {
			out = append(out, storageDeal)
		}
	}

	return out, nil
}

func (n *ClientNodeAdapter) ListStorageProviders(ctx context.Context, encodedTs shared.TipSetToken) ([]*storagemarket.StorageProviderInfo, error) {
	tsk, err := types.TipSetKeyFromBytes(encodedTs)
	if err != nil {
		return nil, err
	}

	addresses, err := n.StateListMiners(ctx, tsk)
	if err != nil {
		return nil, err
	}

	var out []*storagemarket.StorageProviderInfo

	for _, addr := range addresses {
		mi, err := n.GetMinerInfo(ctx, addr, encodedTs)
		if err != nil {
			return nil, err
		}

		out = append(out, mi)
	}

	return out, nil
}

// ValidatePublishedDeal validates that the provided deal has appeared on chain and references the same ClientDeal
// returns the Deal id if there is no error
func (n *ClientNodeAdapter) ValidatePublishedDeal(ctx context.Context, deal storagemarket.ClientDeal) (abi.DealID, error) {
	glog.Info("DEAL ACCEPTED!")

	pubmsg, err := n.cs.GetMessage(*deal.PublishMessage)
	if err != nil {
		return 0, xerrors.Errorf("getting deal pubsish message: %w", err)
	}

	mi, err := nodeapi.StateMinerInfo(ctx, n.node, n.cs.GetHeaviestTipSet(), deal.Proposal.Provider)
	if err != nil {
		return 0, xerrors.Errorf("getting miner worker failed: %w", err)
	}

	fromid, err := n.StateLookupID(ctx, pubmsg.From, types.EmptyTSK)
	if err != nil {
		return 0, xerrors.Errorf("failed to resolve from msg ID addr: %w", err)
	}

	if fromid != mi.Worker {
		return 0, xerrors.Errorf("deal wasn't published by storage provider: from=%s, provider=%s", pubmsg.From, deal.Proposal.Provider)
	}

	if pubmsg.To != builtin.StorageMarketActorAddr {
		return 0, xerrors.Errorf("deal publish message wasn't set to StorageMarket actor (to=%s)", pubmsg.To)
	}

	if pubmsg.Method != builtin.MethodsMarket.PublishStorageDeals {
		return 0, xerrors.Errorf("deal publish message called incorrect method (method=%s)", pubmsg.Method)
	}

	var params samarket.PublishStorageDealsParams
	if err := params.UnmarshalCBOR(bytes.NewReader(pubmsg.Params)); err != nil {
		return 0, err
	}

	dealIdx := -1
	for i, storageDeal := range params.Deals {
		// TODO: make it less hacky
		sd := storageDeal
		eq, err := cborutil.Equals(&deal.ClientDealProposal, &sd)
		if err != nil {
			return 0, err
		}
		if eq {
			dealIdx = i
			break
		}
	}

	if dealIdx == -1 {
		return 0, xerrors.Errorf("deal publish didn't contain our deal (message cid: %s)", deal.PublishMessage)
	}

	// TODO: timeout
	_, ret, err := n.sm.WaitForMessage(ctx, *deal.PublishMessage, build.MessageConfidence)
	if err != nil {
		return 0, xerrors.Errorf("waiting for deal publish message: %w", err)
	}
	if ret.ExitCode != 0 {
		return 0, xerrors.Errorf("deal publish failed: exit=%d", ret.ExitCode)
	}

	var res samarket.PublishStorageDealsReturn
	if err := res.UnmarshalCBOR(bytes.NewReader(ret.Return)); err != nil {
		return 0, err
	}

	return res.IDs[dealIdx], nil
}

func (n *ClientNodeAdapter) SignProposal(ctx context.Context, signer address.Address, proposal market.DealProposal) (*market.ClientDealProposal, error) {
	// TODO: output spec signed proposal
	buf, err := cborutil.Dump(&proposal)
	if err != nil {
		return nil, err
	}

	signer, err = n.StateAccountKey(ctx, signer, types.EmptyTSK)
	if err != nil {
		return nil, err
	}

	sig, err := n.Wallet.Sign(ctx, signer, buf)
	if err != nil {
		return nil, err
	}

	return &samarket.ClientDealProposal{
		Proposal:        proposal,
		ClientSignature: *sig,
	}, nil
}

func (n *ClientNodeAdapter) GetDefaultWalletAddress(ctx context.Context) (address.Address, error) {
	return n.Wallet.GetDefault()
}

func (n *ClientNodeAdapter) ValidateAskSignature(ctx context.Context, ask *storagemarket.SignedStorageAsk, encodedTs shared.TipSetToken) (bool, error) {
	tsk, err := types.TipSetKeyFromBytes(encodedTs)
	if err != nil {
		return false, err
	}

	mi, err := n.StateMinerInfo(ctx, ask.Ask.Miner, tsk)
	if err != nil {
		return false, xerrors.Errorf("failed to get worker for miner in ask %v", err)
	}

	sigb, err := cborutil.Dump(ask.Ask)
	if err != nil {
		return false, xerrors.Errorf("failed to re-serialize ask")
	}

	ts, err := n.ChainGetTipSet(ctx, tsk)
	if err != nil {
		return false, xerrors.Errorf("failed to load tipset")
	}

	m, err := n.sm.ResolveToKeyAddress(ctx, mi.Worker, ts)

	if err != nil {
		return false, xerrors.Errorf("failed to resolve miner to key address")
	}

	err = sigs.Verify(ask.Signature, m, sigb)
	return err == nil, err
}

func (n *ClientNodeAdapter) GetMinerInfo(ctx context.Context, maddr address.Address, encodedTs shared.TipSetToken) (*storagemarket.StorageProviderInfo, error) {
	tsk, err := types.TipSetKeyFromBytes(encodedTs)
	if err != nil {
		return nil, err
	}
	mi, err := n.StateMinerInfo(ctx, maddr, tsk)
	if err != nil {
		return nil, err
	}

	out := utils.NewStorageProviderInfo(maddr, mi.Worker, mi.SectorSize, mi.PeerId, mi.Multiaddrs)
	return &out, nil
}

// GetChainHead returns a tipset token for the current chain head
func (n *ClientNodeAdapter) GetChainHead(ctx context.Context) (shared.TipSetToken, abi.ChainEpoch, error) {
	head, err := n.ChainHead(ctx)
	if err != nil {
		return nil, 0, err
	}

	return head.Key().Bytes(), head.Height(), nil
}

// Adds funds with the StorageMinerActor for a storage participant.  Used by both providers and clients.
func (n *ClientNodeAdapter) AddFunds(ctx context.Context, addr address.Address, amount abi.TokenAmount) (cid.Cid, error) {
	// (Provider Node API)
	smsg, err := n.MpoolPushMessage(ctx, &types.Message{
		To:     builtin.StorageMarketActorAddr,
		From:   addr,
		Value:  amount,
		Method: builtin.MethodsMarket.AddBalance,
	})
	if err != nil {
		return cid.Undef, err
	}

	return smsg.Cid(), nil
}

// EnsureFunds ensures that a storage market participant has a certain amount of available funds
// If additional funds are needed, they will be sent from the 'wallet' address, and a cid for the
// corresponding chain message is returned
func (n *ClientNodeAdapter) EnsureFunds(ctx context.Context, addr, wallet address.Address, amount abi.TokenAmount, tok shared.TipSetToken) (cid.Cid, error) {
	return n.EnsureAvailable(ctx, addr, wallet, amount)
}

// GetBalance returns locked/unlocked for a storage participant.  Used by both providers and clients.
func (n *ClientNodeAdapter) GetBalance(ctx context.Context, addr address.Address, encodedTs shared.TipSetToken) (storagemarket.Balance, error) {
	tsk, err := types.TipSetKeyFromBytes(encodedTs)
	if err != nil {
		return storagemarket.Balance{}, err
	}

	bal, err := n.StateMarketBalance(ctx, addr, tsk)
	if err != nil {
		return storagemarket.Balance{}, err
	}

	return utils.ToSharedBalance(bal), nil
}

// VerifySignature verifies a given set of data was signed properly by a given address's private key
func (n *ClientNodeAdapter) VerifySignature(ctx context.Context, signature crypto.Signature, signer address.Address, plaintext []byte, encodedTs shared.TipSetToken) (bool, error) {
	addr, err := n.StateAccountKey(ctx, signer, types.EmptyTSK)
	if err != nil {
		return false, err
	}

	err = sigs.Verify(&signature, addr, plaintext)
	return err == nil, err
}

// WaitForMessage waits until a message appears on chain. If it is already on chain, the callback is called immediately
func (n *ClientNodeAdapter) WaitForMessage(ctx context.Context, mcid cid.Cid, onCompletion func(exitcode.ExitCode, []byte, error) error) error {
	receipt, err := n.StateWaitMsg(ctx, mcid, build.MessageConfidence)
	if err != nil {
		return onCompletion(0, nil, err)
	}
	return onCompletion(receipt.Receipt.ExitCode, receipt.Receipt.Return, nil)
}

// SignsBytes signs the given data with the given address's private key
func (n *ClientNodeAdapter) SignBytes(ctx context.Context, signer address.Address, b []byte) (*crypto.Signature, error) {
	return nil, nil
}

// OnDealSectorCommitted waits for a deal's sector to be sealed and proved, indicating the deal is active
func (n *ClientNodeAdapter) OnDealSectorCommitted(ctx context.Context, provider address.Address, dealID abi.DealID, cb storagemarket.DealSectorCommittedCallback) error {
	checkFunc := func(ts *types.TipSet) (done bool, more bool, err error) {
		sd, err := nodeapi.GetStorageDeal(ctx, n.node, dealID, ts)

		if err != nil {
			// TODO: This may be fine for some errors
			return false, false, xerrors.Errorf("client: failed to look up deal on chain: %w", err)
		}

		if sd.State.SectorStartEpoch > 0 {
			cb(nil)
			return true, false, nil
		}

		return false, true, nil
	}

	called := func(msg *types.Message, rec *types.MessageReceipt, ts *types.TipSet, curH abi.ChainEpoch) (more bool, err error) {
		defer func() {
			if err != nil {
				cb(xerrors.Errorf("handling applied event: %w", err))
			}
		}()

		if msg == nil {
			glog.Errorf("timed out waiting for deal activation... what now?")
			return false, nil
		}

		sd, err := nodeapi.GetStorageDeal(ctx, n.node, dealID, ts)
		if err != nil {
			return false, xerrors.Errorf("failed to look up deal on chain: %w", err)
		}

		if sd.State.SectorStartEpoch < 1 {
			return false, xerrors.Errorf("deal wasn't active: deal=%d, parentState=%s, h=%d", dealID, ts.ParentState(), ts.Height())
		}

		glog.Infof("Storage deal %d activated at epoch %d", dealID, sd.State.SectorStartEpoch)

		cb(nil)

		return false, nil
	}

	revert := func(ctx context.Context, ts *types.TipSet) error {
		glog.Warningf("deal activation reverted; TODO: actually handle this!")
		// TODO: Just go back to DealSealing?
		return nil
	}

	var sectorNumber abi.SectorNumber
	var sectorFound bool
	matchEvent := func(msg *types.Message) (matchOnce bool, matched bool, err error) {
		if msg.To != provider {
			return true, false, nil
		}

		switch msg.Method {
		case builtin.MethodsMiner.PreCommitSector:
			var params miner.SectorPreCommitInfo
			if err := params.UnmarshalCBOR(bytes.NewReader(msg.Params)); err != nil {
				return true, false, xerrors.Errorf("unmarshal pre commit: %w", err)
			}

			for _, did := range params.DealIDs {
				if did == abi.DealID(dealID) {
					sectorNumber = params.SectorNumber
					sectorFound = true
					return true, false, nil
				}
			}

			return true, false, nil
		case builtin.MethodsMiner.ProveCommitSector:
			var params miner.ProveCommitSectorParams
			if err := params.UnmarshalCBOR(bytes.NewReader(msg.Params)); err != nil {
				return true, false, xerrors.Errorf("failed to unmarshal prove commit sector params: %w", err)
			}

			if !sectorFound {
				return true, false, nil
			}

			if params.SectorNumber != sectorNumber {
				return true, false, nil
			}

			return false, true, nil
		default:
			return true, false, nil
		}
	}

	if err := n.ev.Called(checkFunc, called, revert, int(build.MessageConfidence+1), build.SealRandomnessLookbackLimit, matchEvent); err != nil {
		return xerrors.Errorf("failed to set up called handler: %w", err)
	}
	return nil
}

func (n *ClientNodeAdapter) OnDealExpiredOrSlashed(ctx context.Context, dealID abi.DealID, onDealExpired storagemarket.DealExpiredCallback, onDealSlashed storagemarket.DealSlashedCallback) error {
	head, err := n.ChainHead(ctx)
	if err != nil {
		return xerrors.Errorf("client: failed to get chain head: %w", err)
	}

	sd, err := n.StateMarketStorageDeal(ctx, dealID, head.Key())
	if err != nil {
		return xerrors.Errorf("client: failed to look up deal %d on chain: %w", dealID, err)
	}

	// Called immediately to check if the deal has already expired or been slashed
	checkFunc := func(ts *types.TipSet) (done bool, more bool, err error) {
		// Check if the deal has already expired
		if sd.Proposal.EndEpoch <= ts.Height() {
			onDealExpired(nil)
			return true, false, nil
		}

		// If there is no deal assume it's already been slashed
		if sd.State.SectorStartEpoch < 0 {
			onDealSlashed(ts.Height(), nil)
			return true, false, nil
		}

		// No events have occurred yet, so return
		// done: false, more: true (keep listening for events)
		return false, true, nil
	}

	// Called when there was a match against the state change we're looking for
	// and the chain has advanced to the confidence height
	stateChanged := func(ts *types.TipSet, ts2 *types.TipSet, states events.StateChange, h abi.ChainEpoch) (more bool, err error) {
		// Check if the deal has already expired
		if sd.Proposal.EndEpoch <= ts2.Height() {
			onDealExpired(nil)
			return false, nil
		}

		// Timeout waiting for state change
		if states == nil {
			glog.Error("timed out waiting for deal expiry")
			return false, nil
		}

		changedDeals, ok := states.(state.ChangedDeals)
		if !ok {
			panic("Expected state.ChangedDeals")
		}

		deal, ok := changedDeals[dealID]
		if !ok {
			// No change to deal
			return true, nil
		}

		// Deal was slashed
		if deal.To == nil {
			onDealSlashed(ts2.Height(), nil)
			return false, nil
		}

		return true, nil
	}

	// Called when there was a chain reorg and the state change was reverted
	revert := func(ctx context.Context, ts *types.TipSet) error {
		// TODO: Is it ok to just ignore this?
		glog.Warning("deal state reverted; TODO: actually handle this!")
		return nil
	}

	// Watch for state changes to the deal
	preds := state.NewStatePredicates(n)
	dealDiff := preds.OnStorageMarketActorChanged(
		preds.OnDealStateChanged(
			preds.DealStateChangedForIDs([]abi.DealID{dealID})))
	match := func(oldTs, newTs *types.TipSet) (bool, events.StateChange, error) {
		return dealDiff(ctx, oldTs.Key(), newTs.Key())
	}

	// Wait until after the end epoch for the deal and then timeout
	timeout := (sd.Proposal.EndEpoch - head.Height()) + 1
	if err := n.ev.StateChanged(checkFunc, stateChanged, revert, int(build.MessageConfidence)+1, timeout, match); err != nil {
		return xerrors.Errorf("failed to set up state changed handler: %w", err)
	}

	return nil
}
