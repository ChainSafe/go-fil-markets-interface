// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ChainSafe/go-fil-markets-interface/client"
	"github.com/ChainSafe/go-fil-markets-interface/config"
	"github.com/ChainSafe/go-fil-markets-interface/nodeapi"
	"github.com/ChainSafe/go-fil-markets-interface/utils"
	"github.com/filecoin-project/go-address"
	storagemarket "github.com/filecoin-project/go-fil-markets/storagemarket"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestRetrievalCMD(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)
	err := set.Parse([]string{"QmVMnCY9ic84w7ujYRVANYFH7xnM2YKohMKF66fYA61s2o"})
	require.NoError(t, err)
	fmt.Println(clientFindCmd.Action(ctx))
}

func TestMain(m *testing.M) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config.Load(currentDir + "/../../config/config.json")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMarketStorage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	cctx := cli.NewContext(nil, set, nil)

	nodeapi, nodeCloser, err := nodeapi.GetNodeAPI(cctx)
	if err != nil {
		fmt.Println("Stopping test: Make sure lotus node is up before running this test.")
		return
	}
	defer nodeCloser()

	api, marketCloser, err := client.GetMarketAPI(cctx)
	if err != nil {
		fmt.Println("Stopping test: Make sure go-fil-markets is up before running this test.")
		return
	}
	defer marketCloser()

	mapi, minerCloser, err := client.GetMinerAPI(cctx)
	if err != nil {
		fmt.Println("Stopping test: Make sure lotus miner is up before running this test.")
		return
	}
	defer minerCloser()

	ctx := utils.ReqContext(cctx)

	absPath, err := filepath.Abs("../../data/hello_remote.txt")
	require.NoError(t, err)

	ref := lapi.FileRef{
		Path:  absPath,
		IsCAR: cctx.Bool("car"),
	}
	c, err := api.ClientImport(ctx, ref)
	require.NoError(t, err)

	encoder, err := GetCidEncoder(cctx)
	require.NoError(t, err)

	fmt.Println("Import ", c.ImportID)
	fileCid := encoder.Encode(c.Root)
	fmt.Println("FileCid: ", fileCid)

	fmt.Println("Starting the deal")
	miner, err := address.NewFromString("t01000")
	require.NoError(t, err)

	price, err := types.ParseFIL("0.0000000005")
	require.NoError(t, err)

	dur, err := strconv.ParseInt("518402", 10, 32)
	require.NoError(t, err)

	a, err := nodeapi.WalletDefaultAddress(ctx)
	require.NoError(t, err)

	dataRef := &storagemarket.DataRef{
		TransferType: storagemarket.TTGraphsync,
		Root:         c.Root,
	}
	proposal, err := api.ClientStartDeal(ctx, &lapi.StartDealParams{
		Data:              dataRef,
		Wallet:            a,
		Miner:             miner,
		EpochPrice:        types.BigInt(price),
		MinBlocksDuration: uint64(dur),
		DealStartEpoch:    abi.ChainEpoch(cctx.Int64("start-epoch")),
		FastRetrieval:     cctx.Bool("fast-retrieval"),
		VerifiedDeal:      false,
	})
	require.NoError(t, err)

	encoder, err = GetCidEncoder(cctx)
	require.NoError(t, err)

	fmt.Println("Deal ID: ", encoder.Encode(*proposal))
	for {
		di, err := api.ClientGetDealInfo(ctx, *proposal)
		require.NoError(t, err)
		if di.State == storagemarket.StorageDealSealing {
			dim, _ := json.MarshalIndent(di, "", "  ")
			fmt.Println("Deal Info: ", string(dim))
			break
		}
		fmt.Printf("Waiting for Deal ID %s to reach sealing state.\n", encoder.Encode(*proposal))
		time.Sleep(10 * time.Second)
	}

	fmt.Println("Sealing the sector")
	time.Sleep(10 * time.Second)
	for {
		sector, err := mapi.SectorsList(ctx)
		require.NoError(t, err)
		sealed := true

		// Poll the sectors and wait till all of them are sealed.
		for _, sNum := range sector {
			sInfo, err := mapi.SectorsStatus(ctx, sNum, false)
			require.NoError(t, err)

			if sInfo.State == lapi.SectorState(sealing.WaitDeals) {
				require.NoError(t, mapi.SectorStartSealing(ctx, sNum))
			}

			if sInfo.State != lapi.SectorState(sealing.Proving) {
				sealed = false
				fmt.Println("Waiting for sector ", sNum)
				time.Sleep(10 * time.Second)
			}
		}

		if sealed {
			break
		}
	}
}
