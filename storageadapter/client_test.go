// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package storageadapter

import (
	"bytes"
	"context"
	"errors"
	"github.com/ipfs/go-datastore/namespace"
	"io"
	"testing"

	dtimpl "github.com/filecoin-project/go-data-transfer/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/transport/graphsync"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/discovery"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/funds"
	"github.com/filecoin-project/go-multistore"
	"github.com/filecoin-project/go-storedcounter"
	bstore "github.com/filecoin-project/lotus/lib/blockstore"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-datastore"
	dss "github.com/ipfs/go-datastore/sync"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	"github.com/ipfs/go-graphsync/network"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/require"
)

func TestNewStorageClientNode(t *testing.T) {
	t.Skip() // Skipping as of now.
	// Either user Real server or Mock server implementation to test this.

	ctx := context.Background()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}
	h, err := libp2p.New(ctx, opts...)
	require.NoError(t, err)

	ds := dss.MutexWrap(datastore.NewMapDatastore())
	bs := bstore.NewBlockstore(namespace.Wrap(ds, datastore.NewKey("blockstore")))
	mds, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)

	scn := NewStorageClientNode()

	clientDealFunds, err := funds.NewDealFunds(ds, datastore.NewKey("storage/client/dealfunds"))
	require.NoError(t, err)

	storedCounter := storedcounter.New(ds, datastore.NewKey("counter"))

	makeLoader := func(bs bstore.Blockstore) ipld.Loader {
		return func(lnk ipld.Link, lnkCtx ipld.LinkContext) (io.Reader, error) {
			c, ok := lnk.(cidlink.Link)
			if !ok {
				return nil, errors.New("incorrect Link Type")
			}
			// read block from one store
			block, err := bs.Get(c.Cid)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(block.RawData()), nil
		}
	}

	makeStorer := func(bs bstore.Blockstore) ipld.Storer {
		return func(lnkCtx ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			var buf bytes.Buffer
			var committer ipld.StoreCommitter = func(lnk ipld.Link) error {
				c, ok := lnk.(cidlink.Link)
				if !ok {
					return errors.New("incorrect Link Type")
				}
				block, err := blocks.NewBlockWithCid(buf.Bytes(), c.Cid)
				if err != nil {
					return err
				}
				return bs.Put(block)
			}
			return &buf, committer, nil
		}
	}

	graphSync := graphsyncimpl.New(ctx, network.NewFromLibp2pHost(h), makeLoader(bs), makeStorer(bs))
	transport := dtgstransport.NewTransport(h.ID(), graphSync)
	dt, err := dtimpl.NewDataTransfer(ds, dtnet.NewFromLibp2pHost(h), transport, storedCounter)
	require.NoError(t, err)

	peerResolver := discovery.NewLocal(ds)

	_, err = StorageClient(h, bs, mds, dt, peerResolver, ds, scn, clientDealFunds)
	require.NoError(t, err)
}
