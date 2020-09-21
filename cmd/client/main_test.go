package main

import (
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
	t.Skip()
	set := flag.NewFlagSet("test", 0)
	cctx := cli.NewContext(nil, set, nil)

	nodeapi, nodeCloser, err := nodeapi.GetNodeAPI(cctx)
	require.NoError(t, err)
	defer nodeCloser()

	api, closer, err := client.GetMarketAPI(cctx)
	require.NoError(t, err)
	defer closer()

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

	mapi, closer, err := client.GetMinerAPI(cctx)
	require.NoError(t, err)
	defer closer()

	for {
		sector, err := mapi.SectorsList(ctx)
		require.NoError(t, err)
		sealed := true

		// Poll the sectors and wait till all of them are sealed.
		for _, sNum := range sector {
			sInfo, err := mapi.SectorsStatus(ctx, sNum, false)
			require.NoError(t, err)
			if sInfo.State != "proving" {
				sealed = false
				time.Sleep(5)
			}
		}

		if sealed {
			break
		}
	}
}
