#!/bin/bash

printf "#1) Submit proposal to create stake Alliance...\n\n"
allianced tx gov submit-legacy-proposal create-alliance stake 0.5 0 1 0 0.1 1s --from=node0 --home ../../.testnets/node0/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --deposit=1000000000stake --chain-id=alliance-testnet-1 -y

PROPOSAL_ID=$(allianced query gov proposals --count-total --output json --home ../../.testnets/node0/allianced  --chain-id=alliance-testnet-1 | jq .proposals[-1].id -r)

printf "\n#2) Vote to pass the proposal...\n\n"
allianced tx gov vote $PROPOSAL_ID yes --from=node0 --home ../../.testnets/node0/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=alliance-testnet-1 -y 
allianced tx gov vote $PROPOSAL_ID yes --from=node1 --home ../../.testnets/node1/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=alliance-testnet-1 -y 
allianced tx gov vote $PROPOSAL_ID yes --from=node2 --home ../../.testnets/node2/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=alliance-testnet-1 -y 
