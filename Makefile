#!/usr/bin/make -f

###############################################################################
###                                Localnet                                 ###
###############################################################################

localnet-alliance-rmi:
	docker rmi terra-money/localnet-allianced 2>/dev/null; true

localnet-build-env: localnet-alliance-rmi
	docker build --tag terra-money/localnet-allianced -f scripts/BuildDockerfile \
    		$(shell git rev-parse --show-toplevel)
	

localnet-start-nodes:
	docker run --rm -v $(CURDIR)/.testnet:/data terra-money/localnet-allianced \
			  testnet init-files --v 4 -o /data --starting-ip-address 192.168.10.2 --keyring-backend=test
	docker-compose up -d

localnet-stop:
	docker-compose down

localnet-start: localnet-stop localnet-build-env localnet-build-nodes

.PHONY: localnet-start localnet-stop localnet-build-env localnet-build-nodes
