// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package nodeapi

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/multiformats/go-multiaddr"
)

var (
	fullnodeAPIInfo = "FULLNODE_API_INFO"
	storageAPIInfo  = "STORAGE_API_INFO"
)

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(a.Token))
		return headers
	}
	return http.Header{}
}

func GetAPIInfo() (APIInfo, error) {
	if val, ok := os.LookupEnv(fullnodeAPIInfo); ok {
		return parseEnv(val, fullnodeAPIInfo)
	} else if val, ok := os.LookupEnv(storageAPIInfo); ok {
		return parseEnv(val, storageAPIInfo)
	}

	return APIInfo{}, fmt.Errorf("neither %s or %s are set", fullnodeAPIInfo, storageAPIInfo)
}

func parseEnv(val, env string) (APIInfo, error) {
	sp := strings.SplitN(val, ":", 2)
	if len(sp) != 2 {
		return APIInfo{}, fmt.Errorf("could not parse env(%s)", env)
	}

	ma, err := multiaddr.NewMultiaddr(sp[1])
	if err != nil {
		return APIInfo{}, fmt.Errorf("could not parse multiaddr from env(%s): %w", env, err)
	}
	return APIInfo{
		Addr:  ma,
		Token: []byte(sp[0]),
	}, nil
}
