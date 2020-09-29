PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint
LOTUS_DIR=extern/lotus
FFI_DIR=extern/filecoin-ffi

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
	git -C $(FFI_DIR) reset HEAD --hard
	git -C $(FFI_DIR) checkout cddc566
	git -C $(FFI_DIR) clean -fdx
	make -C $(FFI_DIR)

get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: get-lint submodule
	./bin/golangci-lint run ./... --timeout 5m0s -v --new-from-rev origin/main

test: submodule
	go test ./...

license:
	./scripts/add_license.sh
