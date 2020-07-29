PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

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
	make -C $(FFI_PATH)

$(GOLANGCI):
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: submodule $(GOLANGCI)
	./bin/golangci-lint run ./... --timeout 5m0s

test: submodule
	go test ./...

license:
	./scripts/add_license.sh
