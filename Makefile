#!/usr/bin/make -f
DOCKER := $(shell which docker)
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
PWD := $(shell pwd)
BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

export GO111MODULE = on

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
empty = $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(empty),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=alliance \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=allianced \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/terra-money/alliance/app.Bech32Prefix=alliance \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)' -trimpath

# The below include contains the tools and runsim targets.
#include contrib/devtools/Makefile
all: install lint test

###############################################################################
###                           Build and Install Binary                      ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/ ./cmd/allianced
endif

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/allianced

###############################################################################
###                                Test                                     ###
###############################################################################

test: test-unit test-e2e

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' github.com/terra-money/alliance/x/alliance/keeper/tests github.com/terra-money/alliance/x/alliance/types/tests

test-e2e:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' github.com/terra-money/alliance/x/alliance/tests/e2e

test-benchmark:
	@VERSION=$(VERSION) go test -v -mod=readonly -tags='ledger test_ledger_mock' github.com/terra-money/alliance/x/alliance/tests/benchmark

test-simulate:
	@VERSION=$(VERSION) go test -v -run=TestFullAppSimulation ./app -NumBlocks 200 -BlockSize 10 -Commit -Enabled -Period 1

.PHONY: test test-unit test-e2e test-benchmark test-simulate
###############################################################################
###                                Linting                                  ###
###############################################################################

format-tools:
	go install mvdan.cc/gofumpt@v0.4.0
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

lint: format-tools
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" -not -path "*pb.gw.go" | xargs gofumpt -d

lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.52.2-alpine golangci-lint run --timeout 10m

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" -not -path "*pb.gw.go" | xargs gofumpt -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" -not -path "*pb.gw.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" -not -path "*pb.gw.go" | xargs goimports -w -local github.com/terra-money/alliance

###############################################################################
###                                Protobuf                                 ###
###############################################################################
PROTO_BUILDER_IMAGE=ghcr.io/cosmos/proto-builder

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE) sh ./scripts/protocgen.sh

###############################################################################
###                                Local Testnet (docker)                   ###
###############################################################################

start: install 
	./scripts/local/start.sh

localnet-alliance-rmi:
	$(DOCKER) rmi terra-money/localnet-alliance 2>/dev/null; true

localnet-build-env: localnet-alliance-rmi
	$(DOCKER) build --tag terra-money/localnet-alliance -f scripts/containers/Dockerfile \
    		$(shell git rev-parse --show-toplevel)

localnet-build-nodes:
	$(DOCKER) run --rm -v $(CURDIR)/.testnets:/alliance terra-money/localnet-alliance \
		testnet init-files --v 3 -o /alliance --starting-ip-address 192.168.5.20 --keyring-backend=test --chain-id=alliance-testnet-1
	docker-compose up -d

localnet-stop:
	docker-compose down

localnet-start: localnet-stop localnet-build-env localnet-build-nodes

.PHONY: localnet-start localnet-stop localnet-build-env localnet-build-nodes