#!/bin/bash
COIN_DENOM=stake
VAL_WALLET_ADDRESS=$(allianced --home ../../.testnets/node0/allianced keys show node0 --keyring-backend test -a)
ORIGIN_VAL=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
DEST_VAL=$(allianced query staking validators --output json | jq .validators[1].operator_address --raw-output)

printf "#1) ClaimRewards 10000000000$COIN_DENOM thru x/alliance $COIN_DENOM...\n\n"
allianced tx alliance redelegate $ORIGIN_VAL $DEST_VAL 1000000$COIN_DENOM --from=node0 --home ../../.testnets/node0/allianced --keyring-backend=test --broadcast-mode=block --gas 1000000 --chain-id=alliance-testnet-1 -y

printf "\n#2) Query delegation on x/alliance by delegator, validator and $COIN_DENOM...\n\n"
allianced query alliance delegation $VAL_WALLET_ADDRESS $DEST_VAL $COIN_DENOM
