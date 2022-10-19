#!/bin/bash

CHAIN_DIR=./data
CHAINID=${CHAINID:-alliance}

VAL_WALLET_ADDRESS=$(allianced --home $CHAIN_DIR/$CHAINID keys show val1 --keyring-backend test -a)
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)

printf "\n\n#1)Query wallet balances...\n\n"
allianced query bank balances $VAL_WALLET_ADDRESS --home $CHAIN_DIR/$CHAINID

printf "\n\n#2)Query rewards x/alliance...\n\n"
allianced query alliance rewards $VAL_WALLET_ADDRESS $VAL_ADDR token --home $CHAIN_DIR/$CHAINID

printf "\n\n#3)Claim rewards x/alliance...\n\n"
allianced tx alliance claim-rewards $VAL_ADDR token --from=val1 --home $CHAIN_DIR/$CHAINID --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n\n#4)Query rewards x/alliance...\n\n"
allianced query alliance rewards $VAL_WALLET_ADDRESS $VAL_ADDR token --home $CHAIN_DIR/$CHAINID

printf "\n\n#5)Query wallet balances after claim...\n\n"
allianced query bank balances $VAL_WALLET_ADDRESS --home $CHAIN_DIR/$CHAINID