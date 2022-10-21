#!/bin/bash

VAL_WALLET_ADDRESS=$(allianced --home ./data/alliance keys show val1 --keyring-backend test -a)
DEMO_WALLET_ADDRESS=$(allianced --home ./data/alliance keys show demowallet1 --keyring-backend test -a)

VAL_ADDR=$(allianced query staking validators --output json | jq .validators[0].operator_address --raw-output)
TOKEN_DENOM=ulunax

printf "\n\n#1) Query wallet balances...\n\n"
allianced query bank balances $DEMO_WALLET_ADDRESS --home ./data/alliance

#printf "\n\n#2) Query rewards x/alliance...\n\n"
#allianced query alliance rewards $DEMO_WALLET_ADDRESS $VAL_ADDR $TOKEN_DENOM --home ./data/alliance

#printf "\n\n#3) Query native staked rewards...\n\n"
#allianced query distribution rewards $DEMO_WALLET_ADDRESS $VAL_ADDR --home ./data/alliance

#printf "\n\n#4) Claim rewards from validator...\n\n"
#allianced tx distribution withdraw-rewards $VAL_ADDR --from=demowallet1 --home ./data/alliance --keyring-backend=test --broadcast-mode=block --gas 1000000 -y #> /dev/null 2>&1

printf "\n\n#2) Claim rewards from x/alliance $TOKEN_DENOM...\n\n"
allianced tx alliance claim-rewards $VAL_ADDR $TOKEN_DENOM --from=demowallet1 --home ./data/alliance --keyring-backend=test --broadcast-mode=block --gas 1000000 -y > /dev/null 2>&1

#printf "\n\n#6) Query rewards x/alliance...\n\n"
#allianced query alliance rewards $DEMO_WALLET_ADDRESS $VAL_ADDR $TOKEN_DENOM --home ./data/alliance

#printf "\n\n#7) Query native staked rewards...\n\n"
#allianced query distribution rewards $DEMO_WALLET_ADDRESS $VAL_ADDR --home ./data/alliance

printf "\n\n#3) Query wallet balances after claim...\n\n"
allianced query bank balances $DEMO_WALLET_ADDRESS --home ./data/alliance
