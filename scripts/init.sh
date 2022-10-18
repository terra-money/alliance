#!/bin/bash

BINARY=${BINARY:-allianced}
CHAIN_DIR=./data
CHAINID=${CHAINID:-test-1}

# alliancevaloper1phaxpevm5wecex2jyaqty2a4v02qj7qmut9cku / alliance1phaxpevm5wecex2jyaqty2a4v02qj7qm24tyvq
VAL_MNEMONIC_1="satisfy adjust timber high purchase tuition stool faith fine install that you unaware feed domain license impose boss human eager hat rent enjoy dawn"

# alliance1cyyzpxplxdzkeea7kwsydadg87357qnaznl0wd
DEMO_MNEMONIC_1="notice oak worry limit wrap speak medal online prefer cluster roof addict wrist behave treat actual wasp year salad speed social layer crew genius"
DEMO_MNEMONIC_2="quality vacuum heart guard buzz spike sight swarm shove special gym robust assume sudden deposit grid alcohol choice devote leader tilt noodle tide penalty"
DEMO_MNEMONIC_3="symbol force gallery make bulk round subway violin worry mixture penalty kingdom boring survey tool fringe patrol sausage hard admit remember broken alien absorb"
DEMO_MNEMONIC_4="bounce success option birth apple portion aunt rural episode solution hockey pencil lend session cause hedgehog slender journey system canvas decorate razor catch empty"

STAKEDENOM=${STAKEDENOM:-stake}
UNBONDING_TIME="60s"
INFLATION="0.950000000000000000"
GOV_PERIOD="60s"

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $CHAIN_DIR/$CHAINID &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $CHAIN_DIR/$CHAINID 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID..."
$BINARY init test --home $CHAIN_DIR/$CHAINID --chain-id=$CHAINID

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $BINARY keys add val1 --home $CHAIN_DIR/$CHAINID --recover --keyring-backend=test
echo $DEMO_MNEMONIC_1 | $BINARY keys add demowallet1 --home $CHAIN_DIR/$CHAINID --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $BINARY keys add demowallet2 --home $CHAIN_DIR/$CHAINID --recover --keyring-backend=test
echo $DEMO_MNEMONIC_3 | $BINARY keys add demowallet3 --home $CHAIN_DIR/$CHAINID --recover --keyring-backend=test
echo $DEMO_MNEMONIC_4 | $BINARY keys add demowallet4 --home $CHAIN_DIR/$CHAINID --recover --keyring-backend=test

$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID keys show val1 --keyring-backend test -a) 100000000000${STAKEDENOM}  --home $CHAIN_DIR/$CHAINID
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID keys show demowallet1 --keyring-backend test -a) 100000000000${STAKEDENOM}  --home $CHAIN_DIR/$CHAINID
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID keys show demowallet2 --keyring-backend test -a) 100000000000${STAKEDENOM}  --home $CHAIN_DIR/$CHAINID
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID keys show demowallet3 --keyring-backend test -a) 100000000000${STAKEDENOM}  --home $CHAIN_DIR/$CHAINID
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID keys show demowallet4 --keyring-backend test -a) 100000000000${STAKEDENOM}  --home $CHAIN_DIR/$CHAINID

echo "Creating and collecting gentx..."
$BINARY gentx val1 7000000000${STAKEDENOM} --home $CHAIN_DIR/$CHAINID --chain-id $CHAINID --keyring-backend test
$BINARY collect-gentxs --home $CHAIN_DIR/$CHAINID

sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $CHAIN_DIR/$CHAINID/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID/config/app.toml
sed -i -e "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.0025$STAKEDENOM\"/g" $CHAIN_DIR/$CHAINID/config/app.toml
sed -i -e 's/enabled = false/enabled = true/g' $CHAIN_DIR/$CHAINID/config/app.toml
sed -i -e 's/prometheus-retention-time = 0/prometheus-retention-time = 1000/g' $CHAIN_DIR/$CHAINID/config/app.toml

## DENOMS
sed -i -e "s/\"denom\": \"stake\",/\"denom\": \"$STAKEDENOM\",/g" $CHAIN_DIR/$CHAINID/config/genesis.json
sed -i -e "s/\"mint_denom\": \"stake\",/\"mint_denom\": \"$STAKEDENOM\",/g" $CHAIN_DIR/$CHAINID/config/genesis.json
sed -i -e "s/\"bond_denom\": \"stake\"/\"bond_denom\": \"$STAKEDENOM\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json

## MINT
sed -i -e "s/\"inflation\": \"0.130000000000000000\"/\"inflation\": \"$INFLATION\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json
sed -i -e "s/\"inflation_rate_change\": \"0.130000000000000000\"/\"inflation_rate_change\": \"$INFLATION\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json
sed -i -e "s/\"inflation_max\": \"0.200000000000000000\"/\"inflation_max\": \"$INFLATION\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json

## STAKING
sed -i -e "s/\"unbonding_time\": \"1814400s\"/\"unbonding_time\": \"$UNBONDING_TIME\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json

## GOV
sed -i -e "s/\"max_deposit_period\": \"172800s\"/\"max_deposit_period\": \"$GOV_PERIOD\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json
sed -i -e "s/\"voting_period\": \"172800s\"/\"voting_period\": \"$GOV_PERIOD\"/g" $CHAIN_DIR/$CHAINID/config/genesis.json