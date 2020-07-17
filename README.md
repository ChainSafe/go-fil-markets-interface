# Filecoin Markets Interface

[![Build Status](https://travis-ci.com/ChainSafe/fil-markets-interface.svg?token=tppYFL7pBXTmrb45py5Q&branch=main)](https://travis-ci.com/ChainSafe/fil-markets-interface)

This project provides an interface for [go-fil-markets](https://github.com/filecoin-project/go-fil-markets) to enable compatibility with any Filecoin node implementation.


## Background

Presently go-fil-markets uses Lotus directly. The interfaces go-fil-markets depends on are well-defined, thus to enable other nodes implementations this project provides a JSON-RPC client that implements the [StorageClientNode](https://github.com/filecoin-project/go-fil-markets/blob/e4e257f097707e6f93f4bf0b5f9aa931871c94a7/storagemarket/nodes.go#L84-L107) and [RetrievalClientNode](https://github.com/filecoin-project/go-fil-markets/blob/e4e257f097707e6f93f4bf0b5f9aa931871c94a7/retrievalmarket/nodes.go#L15-L41) interfaces. This allows the Storage Market Client and the Retrieval Market Client to run independantly of the node process.


### Compatibility 

Compatible nodes will need to support a sub-set of the full node API. They also need to support token authentication for requests.

[TODO: Add specific API requirements]

## Install

## Usage

## License

This repo is dual licensed under [MIT](/LICENSE-MIT) and [Apache 2.0](/LICENSE-APACHE).
