PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint
LOTUS_DIR=extern/lotus

.PHONY: help lint test
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

FFI_PATH:=./extern/filecoin-ffi/

submodule:
	git submodule update --init --recursive
	make -C extern/filecoin-ffi
	git -C $(LOTUS_DIR) reset HEAD --hard
	git -C $(LOTUS_DIR) checkout v0.5.4
	make -C $(LOTUS_DIR) clean
	make -C $(LOTUS_DIR) 2k

get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: get-lint submodule
	./bin/golangci-lint run ./... --timeout 5m0s -v --new-from-rev origin/main

test: submodule
	go test ./...

storagetest: submodule
	./run_lotus.sh &
	./run_lotus_miner.sh &

license:
	./scripts/add_license.sh
