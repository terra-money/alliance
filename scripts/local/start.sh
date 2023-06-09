#!/bin/bash

BINARY=allianced
CHAIN_DIR=$(pwd)/.testnet
CHAINID_1=alliance-testnet-1

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
WALLET_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"

P2PPORT_1=16656
RPCPORT_1=16657
RESTPORT_1=1316
ROSETTA_1=8080
GRPCPORT_1=8090
GRPCWEB_1=8091

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $CHAIN_DIR/$CHAINID_1 

if ! mkdir -p $CHAIN_DIR/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

$BINARY init test --home $CHAIN_DIR/$CHAINID_1 --chain-id=$CHAINID_1 

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $BINARY keys add val1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_1 | $BINARY keys add wallet1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test

VAL1_ADDR=$($BINARY keys show val1 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)
WALLET1_ADDR=$($BINARY keys show wallet1 --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test -a)

$BINARY add-genesis-account $VAL1_ADDR 1000000000000stake --home $CHAIN_DIR/$CHAINID_1
$BINARY add-genesis-account $WALLET1_ADDR 1000000000000stake --home $CHAIN_DIR/$CHAINID_1

echo "Creating and collecting gentx..."
$BINARY gentx val1 6500000000stake --home $CHAIN_DIR/$CHAINID_1 --chain-id $CHAINID_1 --keyring-backend test
sleep 2
$BINARY collect-gentxs --home $CHAIN_DIR/$CHAINID_1 

echo "Changing app.toml defaults and ports in and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's#"tcp://localhost:1317"#"tcp://0.0.0.0:'"$RESTPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/app.toml

echo "Changing genesis.json defaults..."
sed -i -e 's/"voting_period": "172800s"/"voting_period": "10s"/g' $CHAIN_DIR/$CHAINID_1/config/genesis.json
sed -i -e 's/"reward_delay_time": "604800s"/"reward_delay_time": "10s"/g' $CHAIN_DIR/$CHAINID_1/config/genesis.json

echo "Starting $CHAINID_1 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_1.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID_1 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT_1" --grpc-web.address="0.0.0.0:$GRPCWEB_1" > $CHAIN_DIR/$CHAINID_1.log 2>&1 &

