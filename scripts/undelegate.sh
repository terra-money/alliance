#!/bin/bash

DEMO_WALLET_ADDRESS=$(allianced --home ./data/alliance keys show demowallet1 --keyring-backend test -a)
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
COIN_DENOM=ulunax
COIN_AMOUNT=$(allianced query alliance delegation $DEMO_WALLET_ADDRESS $VAL_ADDR $COIN_DENOM --home ./data/alliance --output json | jq .delegation.balance.amount --raw-output | sed 's/\.[0-9]*//')
COINS=5000000000$COIN_DENOM

# FIX: failed to execute message; message index: 0: invalid shares amount: invalid
printf "#1) Undelegate 5000000000$COIN_DENOM from x/alliance $COIN_DENOM...\n\n"
allianced tx alliance undelegate $VAL_ADDR $COINS --from=demowallet1 --home ./data/alliance --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n#2) Query delegations from x/alliance $COIN_DENOM...\n\n"
allianced query alliance alliance $COIN_DENOM

printf "\n#3) Query delegation on x/alliance by delegator, validator and $COIN_DENOM...\n\n"
allianced query alliance delegation $DEMO_WALLET_ADDRESS $VAL_ADDR $COIN_DENOM --home ./data/alliance
