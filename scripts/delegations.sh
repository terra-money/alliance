#!/bin/bash

CHAIN_DIR=./data
CHAINID=${CHAINID:-alliance}

VAL_WALLET_ADDRESS=$(allianced --home $CHAIN_DIR/$CHAINID keys show val1 --keyring-backend test -a)
VAL_ADRESS=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)

echo $VAL_WALLET_ADDRESS
echo $VAL_ADRESS

printf "#1)Delegate thru x/alliance...\n\n"
allianced tx alliance delegate $VAL_ADRESS 10000000stake --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block -y

printf "\n#2)Query delegations from x/alliance by alliance...\n\n"
allianced query alliance alliance stake

printf "\n#4)Query delegations x/alliance by delegator address...\n\n"
allianced query alliance delegationsByDelegator $VAL_WALLET_ADDRESS

printf "\n#5)Query delegations x/alliance by delegator address and validator...\n\n"
allianced query alliance delegationsByDelegatorAndValidator $VAL_WALLET_ADDRESS $VAL_ADRESS

printf "\n#6)Query delegation on x/alliance by delegator, validator and stake...\n\n"
allianced query alliance delegation $VAL_WALLET_ADDRESS $VAL_ADRESS stake
