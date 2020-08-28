// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"context"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

func GetStorageDeal(ctx context.Context, client Node, dealID abi.DealID, ts *types.TipSet) (*api.MarketDeal, error) {
	return NodeClient.UtilsAPI.StateMarketStorageDeal(ctx, dealID, ts.Key())
}

func StateMinerInfo(ctx context.Context, client Node, ts *types.TipSet, maddr address.Address) (api.MinerInfo, error) {
	return NodeClient.UtilsAPI.StateMinerInfo(ctx, maddr, ts.Key())
}

func NewStorageProviderInfo(address address.Address, miner address.Address, sectorSize abi.SectorSize, peer peer.ID, addrs []abi.Multiaddrs) storagemarket.StorageProviderInfo {
	multiaddrs := make([]multiaddr.Multiaddr, 0, len(addrs))
	if addrs == nil {
		peerInfo, _ := NodeClient.UtilsAPI.NetFindPeer(context.Background(), peer)
		multiaddrs = peerInfo.Addrs
	} else {
		for _, a := range addrs {
			maddr, err := multiaddr.NewMultiaddrBytes(a)
			if err != nil {
				return storagemarket.StorageProviderInfo{}
			}
			multiaddrs = append(multiaddrs, maddr)
		}
	}

	return storagemarket.StorageProviderInfo{
		Address:    address,
		Worker:     miner,
		SectorSize: uint64(sectorSize),
		PeerID:     peer,
		Addrs:      multiaddrs,
	}
}

func NewLibP2PHost() (host.Host, error) {
	ctx := context.Background()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}
	return libp2p.New(ctx, opts...)
}
