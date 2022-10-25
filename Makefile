#!/usr/bin/make -f

localnet-alliance-rmi:
	docker rmi terra-money/localnet-alliance 2>/dev/null; true

localnet-build-env: localnet-alliance-rmi
	docker build --tag terra-money/localnet-alliance -f scripts/containers/Dockerfile \
    		$(shell git rev-parse --show-toplevel)
	
localnet-build-nodes:
	docker run --rm -v $(CURDIR)/.testnets:/alliance terra-money/localnet-alliance \
		testnet init-files --v 4 -o /alliance --starting-ip-address 192.168.5.20 --keyring-backend=test
	docker-compose up -d

localnet-stop:
	docker-compose down

localnet-start: localnet-stop localnet-build-env localnet-build-nodes

.PHONY: localnet-start localnet-stop localnet-build-env localnet-build-nodes
