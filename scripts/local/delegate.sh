#!/bin/bash

COIN_DENOM=stake
VAL_WALLET_ADDRESS=$(allianced --home ../../.testnets/node0/allianced keys show node0 --keyring-backend test -a)
VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)

printf "#1) Delegate 10000000000$COIN_DENOM thru x/alliance $COIN_DENOM...\n\n"
allianced tx alliance delegate $VAL_ADDR 10000000000$COIN_DENOM --from=node0 --home ../../.testnets/node0/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=alliance-testnet-1 -y

printf "\n#2) Query delegation on x/alliance by delegator, validator and $COIN_DENOM...\n\n"
allianced query alliance delegation $VAL_WALLET_ADDRESS $VAL_ADDR $COIN_DENOM
