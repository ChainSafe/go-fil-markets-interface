package config

import (
	"encoding/json"
	"fmt"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr/net"
	"io/ioutil"
	"net/http"
	os "os"
	"strings"
)

var (
	nodeAPIInfo   = "NODE_API_INFO"
	marketAPIInfo = "MARKET_API_INFO"
	minerAPIInfo  = "MINER_API_INFO"
	tokenKey      = "Token"
	AddrKey       = "Addr"
	Api           API // TODO(arijit): Remove it from global vars.
)

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
}

type API struct {
	Node   *APIInfo
	Market *APIInfo
	Miner  *APIInfo
}

func (a *APIInfo) UnmarshalJSON(data []byte) error {
	am := make(map[string]interface{})
	if err := json.Unmarshal(data, &am); err != nil {
		return err
	}
	if token, ok := am[tokenKey]; ok {
		a.Token = []byte(token.(string))
	}

	addr, ok := am[AddrKey]
	if !ok {
		return nil
	}

	var err error
	if a.Addr, err = multiaddr.NewMultiaddr(addr.(string)); err != nil {
		return err
	}
	return nil
}

func (a APIInfo) DialArgs() (string, error) {
	_, addr, err := manet.DialArgs(a.Addr)
	return "ws://" + addr + "/rpc/v0", err
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(a.Token))
		return headers
	}
	return http.Header{}
}

func GetAPIInfo() (API, error) {
	var api API
	var err error
	if val, ok := os.LookupEnv(nodeAPIInfo); ok {
		if api.Node, err = parseEnv(val, nodeAPIInfo); err != nil {
			return API{}, err
		}
	}

	if val, ok := os.LookupEnv(marketAPIInfo); ok {
		if api.Market, err = parseEnv(val, marketAPIInfo); err != nil {
			return API{}, err
		}
	}

	if val, ok := os.LookupEnv(minerAPIInfo); ok {
		if api.Miner, err = parseEnv(val, minerAPIInfo); err != nil {
			return API{}, err
		}
	}

	if api.Market == nil || api.Node == nil || api.Miner == nil {
		return API{}, fmt.Errorf("either %s or %s are set", nodeAPIInfo, marketAPIInfo)
	}
	return api, nil
}

func parseEnv(val, env string) (*APIInfo, error) {
	sp := strings.SplitN(val, ":", 2)
	if len(sp) != 2 {
		return &APIInfo{}, fmt.Errorf("could not parse env(%s)", env)
	}

	ma, err := multiaddr.NewMultiaddr(sp[1])
	if err != nil {
		return &APIInfo{}, fmt.Errorf("could not parse multiaddr from env(%s): %w", env, err)
	}
	return &APIInfo{
		Addr:  ma,
		Token: []byte(sp[0]),
	}, nil
}

func Load(configFile string) {
	var err error
	Api, err = GetAPIInfo()
	if err == nil {
		return
	}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to load config from file : %v", err))
	}

	err = json.Unmarshal(file, &Api)
	if err != nil {
		panic(fmt.Sprintf("Unable to unmarshall config from file : %v", err))
	}
}

func (a *API) Store(configFile string) {
	file, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		panic(fmt.Sprintf("Unable to marshall config : %v", err))
	}

	err = ioutil.WriteFile(configFile, file, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to store config to file : %v", err))
	}
}
