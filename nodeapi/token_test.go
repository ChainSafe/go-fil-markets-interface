// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"os"
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestGetAPIInfo(t *testing.T) {
	testAddr := "/ip4/192.168.0.102/tcp/39393/p2p/12D3KooWAvVVGQU8KyTMCvsGoEVnQyztCZWWa2j7HjvbHB7fstjC"
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIl19.jFiihZ1W3W0ttZkXZOzGIr9uW3aq2eqIdEsvdTYrDNQ"
	testVal := testToken + ":" + testAddr
	err := os.Setenv(fullnodeAPIInfo, testVal)
	require.NoError(t, err)

	info, err := GetAPIInfo()
	require.NoError(t, err)

	addr, err := multiaddr.NewMultiaddr(testAddr)
	require.NoError(t, err)

	expected := APIInfo{
		Addr:  addr,
		Token: []byte(testToken),
	}
	require.Equal(t, expected, info)

	err = os.Unsetenv(fullnodeAPIInfo)
	require.NoError(t, err)

	err = os.Setenv(storageAPIInfo, testVal)
	require.NoError(t, err)

	info, err = GetAPIInfo()
	require.NoError(t, err)

	addr, err = multiaddr.NewMultiaddr(testAddr)
	require.NoError(t, err)

	expected = APIInfo{
		Addr:  addr,
		Token: []byte(testToken),
	}
	require.Equal(t, expected, info)
}
