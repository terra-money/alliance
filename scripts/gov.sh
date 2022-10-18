#!/bin/bash
BINARY=${BINARY:-allianced}
CHAINID=${CHAINID:-test-1}
FROM_WALLET=${FROM_WALLET:-val1}

printf "Submit proposal...\n"
$BINARY tx gov submit-proposal gov.json --from=$FROM_WALLET --chain-id=$CHAINID -y

printf "\nDeposit the min funds...\n"
$BINARY tx gov deposit 1 100000000stake --from=$FROM_WALLET --chain-id=$CHAINID -y

printf "\nVote to pass the proposal ...\n"
$BINARY tx gov vote 1 yes --from=$FROM_WALLET --chain-id=$CHAINID -y

printf "\nQuerying alliances...\n"
$BINARY query alliance alliances