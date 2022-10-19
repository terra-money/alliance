#!/bin/bash

VAL_WALLET_ADDRESS=$(allianced --home ./data/alliance keys show val1 --keyring-backend test -a)
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
ALL_TOKENS=$(allianced query alliance delegation $VAL_WALLET_ADDRESS $VAL_ADDR token --home ./data/alliance --output json | jq .delegation.balance.amount --raw-output)
TOKEN=token
TOKENS=$ALL_TOKENS$TOKEN

printf "#1)Undelegate from x/alliance...\n\n"
allianced tx alliance undelegate $VAL_ADDR $TOKENS --from=val1 --home ./data/alliance --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

printf "\n#2)Query delegations from x/alliance by alliance token...\n\n"
allianced query alliance alliance token
