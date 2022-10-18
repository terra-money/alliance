#!/bin/bash
BINARY=${BINARY:-allianced}
CHAINID=${CHAINID:-test-1}

FROM_WALLET=${FROM_WALLET:-demowallet1}
VALIDATOR_ADDR=${VALIDATOR_ADDR:-alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku}

printf "Delegate thru x/alliance..."
$BINARY tx alliance delegate $VALIDATOR_ADDR 100000stake --from=$FROM_WALLET --chain-id=$CHAINID -y

printf "\nCheck the delegations on x/alliance by stake...\n"
$BINARY query alliance alliance stake