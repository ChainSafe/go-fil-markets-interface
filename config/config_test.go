package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	Load("config.json")
	require.NotNil(t, Market.NodeIP)
	require.Equal(t, Market.NodeAuthToken, "")
	require.NotNil(t, Market.MarketIP)
	require.Equal(t, Market.MarketAuthToken, "")
}
