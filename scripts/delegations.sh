#!/bin/bash

printf "#1)Delegate thru x/alliance...\n\n"
allianced tx alliance delegate alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku 1000000stake --from=val1 -y
allianced tx alliance delegate alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku 1000000stake --from=demowallet1 -y

printf "\n#2)Query delegations from x/alliance by alliance...\n\n"
allianced query alliance alliance stake

printf "\n#4)Query delegations x/alliance by delegator address...\n\n"
allianced query alliance delegationsByDelegator alliance1phaxpevm5wecex2jyaqty2a4v02qj7qm24tyvq

printf "\n#5)Query delegations x/alliance by delegator address and validator...\n\n"
allianced query alliance delegationsByDelegatorAndValidator alliance1phaxpevm5wecex2jyaqty2a4v02qj7qm24tyvq alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku

printf "\n#6)Query delegation on x/alliance by delegator, validator and stake...\n\n"
allianced query alliance delegation alliance1phaxpevm5wecex2jyaqty2a4v02qj7qm24tyvq alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku stake