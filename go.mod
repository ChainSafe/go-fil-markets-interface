module github.com/ChainSafe/go-fil-markets-interface

go 1.13

require (
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-cbor-util v0.0.0-20191219014500-08c40a1e63a2
	github.com/filecoin-project/go-fil-markets v0.5.1
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200602181149-522144ab4e24
	github.com/filecoin-project/lotus v0.4.3-0.20200728010654-c0f0e2ba45cf
	github.com/filecoin-project/specs-actors v0.8.1-0.20200724015154-3c690d9b7e1d
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/ipfs/go-cid v0.0.6
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multiaddr-net v0.1.5
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
