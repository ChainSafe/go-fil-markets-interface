// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package config

import (
	"os"
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	Load("config.json")
	require.NotNil(t, Api.Market.Addr)
	require.NotNil(t, Api.Market.Token)
	require.NotNil(t, Api.Node.Addr)
	require.NotNil(t, Api.Node.Token)
}

func TestGetAPIInfo(t *testing.T) {
	testAddr := "/ip4/192.168.0.102/tcp/39393/p2p/12D3KooWAvVVGQU8KyTMCvsGoEVnQyztCZWWa2j7HjvbHB7fstjC"
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIl19.jFiihZ1W3W0ttZkXZOzGIr9uW3aq2eqIdEsvdTYrDNQ"
	testVal := testToken + ":" + testAddr
	err := os.Setenv(nodeAPIInfo, testVal)
	require.NoError(t, err)

	err = os.Setenv(marketAPIInfo, testVal)
	require.NoError(t, err)

	err = os.Setenv(minerAPIInfo, testVal)
	require.NoError(t, err)

	result, err := GetAPIInfo()
	require.NoError(t, err)

	addr, err := multiaddr.NewMultiaddr(testAddr)
	require.NoError(t, err)

	info := &APIInfo{
		Addr:  addr,
		Token: []byte(testToken),
	}

	expected := API{
		Node:   info,
		Market: info,
		Miner:  info,
	}
	require.Equal(t, expected, result)

	err = os.Unsetenv(marketAPIInfo)
	require.NoError(t, err)

	_, err = GetAPIInfo()
	require.Error(t, err)

	err = os.Unsetenv(nodeAPIInfo)
	require.NoError(t, err)

	_, err = GetAPIInfo()
	require.Error(t, err)
}
