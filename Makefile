#!/usr/bin/make -f
DOCKER := $(shell which docker)

###############################################################################
###                                Protobuf                                 ###
###############################################################################
PROTO_BUILDER_IMAGE=tendermintdev/sdk-proto-gen:v0.7

proto-all: proto-swagger-gen proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE) sh ./scripts/protocgen.sh

###############################################################################
###                                Local Testnet (docker)                   ###
###############################################################################

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