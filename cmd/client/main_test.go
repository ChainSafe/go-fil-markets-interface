package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/ChainSafe/go-fil-markets-interface/config"
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
