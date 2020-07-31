module github.com/ChainSafe/go-fil-markets-interface

go 1.13

require (
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-cbor-util v0.0.0-20191219014500-08c40a1e63a2
	github.com/filecoin-project/go-fil-markets v0.5.1
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200602181149-522144ab4e24
	github.com/filecoin-project/go-multistore v0.0.1
	github.com/filecoin-project/lotus v0.4.3-0.20200728010654-c0f0e2ba45cf
	github.com/filecoin-project/sector-storage v0.0.0-20200727112136-9377cb376d25
	github.com/filecoin-project/specs-actors v0.8.1-0.20200724015154-3c690d9b7e1d
	github.com/gbrlsnchs/jwt/v3 v3.0.0-beta.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/ipfs/go-blockservice v0.1.4-0.20200624145336-a978cec6e834
	github.com/ipfs/go-cid v0.0.6
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-ipfs-chunker v0.0.5
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-merkledag v0.3.1
	github.com/ipfs/go-unixfs v0.2.4
	github.com/ipld/go-car v0.1.1-0.20200526133713-1c7508d55aae
	github.com/ipld/go-ipld-prime v0.0.2-0.20200428162820-8b59dc292b8e
	github.com/libp2p/go-libp2p-core v0.6.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multiaddr-net v0.1.5
	github.com/multiformats/go-multihash v0.0.14
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.9.0
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
