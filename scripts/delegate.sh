#!/bin/bash

CHAIN_DIR=./data
CHAINID=${CHAINID:-alliance}

VAL_WALLET_ADDRESS=$(allianced --home $CHAIN_DIR/$CHAINID keys show val1 --keyring-backend test -a)
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)

printf "#1)Delegate thru x/alliance...\n\n"
allianced tx alliance delegate $VAL_ADDR 10000000token --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n#2)Query delegations from x/alliance by alliance...\n\n"
allianced query alliance alliance token

printf "\n#4)Query delegations x/alliance by delegator address...\n\n"
allianced query alliance delegations-by-delegator $VAL_WALLET_ADDRESS

printf "\n#5)Query delegations x/alliance by delegator address and validator...\n\n"
allianced query alliance delegations-by-delegator-and-validator $VAL_WALLET_ADDRESS $VAL_ADDR

printf "\n#6)Query delegation on x/alliance by delegator, validator and token...\n\n"
allianced query alliance delegation $VAL_WALLET_ADDRESS $VAL_ADDR token
