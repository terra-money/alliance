#!/bin/bash

CHAIN_DIR=./data
CHAINID=${CHAINID:-alliance}

printf "#1)Submit proposal...\n\n"
allianced tx gov submit-proposal gov.json --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

sleep 2
PROPOSAL_ID=$(allianced query gov proposals --count-total --output json | jq .pagination.total -r)

printf "\n#2)Deposit the min funds...\n\n"
allianced tx gov deposit $PROPOSAL_ID 10000000stake --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n#3)Vote to pass the proposal ...\n\n"
allianced tx gov vote $PROPOSAL_ID yes --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n#4)Query proposals...\n\n"
allianced query gov proposal $PROPOSAL_ID --home $CHAIN_DIR/$CHAINID

printf "\n#5)Querying alliances...\n\n"
allianced query alliance alliances --home $CHAIN_DIR/$CHAINID

printf "\n#6)Witing for gov proposal to pass...\n\n"
sleep 5

printf "\n#7)Querying alliances after gov passed...\n\n"
allianced query alliance alliances --home $CHAIN_DIR/$CHAINID