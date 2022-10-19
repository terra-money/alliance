#!/bin/bash

printf "#1)Submit proposal...\n\n"
allianced tx gov submit-proposal gov.json --from=val1 -y

printf "\n#2)Deposit the min funds...\n\n"
allianced tx gov deposit 1 10000000stake --from=val1 -y

printf "\n#3)Vote to pass the proposal ...\n\n"
allianced tx gov vote 1 yes --from=val1 -y

printf "\n#4)Query proposals...\n\n"
allianced query gov proposal 1

printf "\n#5)Querying alliances...\n\n"
allianced query alliance alliances

printf "\n#6)Witing for gov proposal to pass...\n\n"
sleep 11

printf "\n#7)Querying alliances after gov passed...\n\n"
allianced query alliance alliances