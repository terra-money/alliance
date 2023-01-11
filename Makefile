#!/usr/bin/make -f

ACC_PREFIX = alliance
CHAIN_ID = alliance-testnet-1
BOND_DENOM = stake

localnet-alliance-rmi:
	docker rmi terra-money/localnet-alliance 2>/dev/null; true

localnet-build-env:
	docker build --tag terra-money/localnet-alliance -f scripts/containers/Dockerfile --build-arg ACC_PREFIX=$(ACC_PREFIX)\
    		$(shell git rev-parse --show-toplevel)

localnet-init: 
	docker run --rm -v $(CURDIR)/.testnets:/alliance terra-money/localnet-alliance \
		testnet init-files --v 3 -o /alliance --starting-ip-address 192.168.5.20 --keyring-backend=test --chain-id=$(CHAIN_ID) --bond-denom=$(BOND_DENOM) --minimum-gas-prices=0$(BOND_DENOM)

localnet-up:
	docker-compose up -d
	
localnet-build-nodes: localnet-init localnet-up

localnet-stop:
	docker-compose down

localnet-start: localnet-stop localnet-build-env localnet-build-nodes

.PHONY: localnet-start localnet-stop localnet-build-env localnet-build-nodes
