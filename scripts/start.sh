#!/bin/bash

BINARY=${BINARY:-allianced}
CHAIN_DIR=./data
CHAINID=${CHAINID:-alliance}
GRPCPORT=9090
GRPCWEB=9091


echo "Starting $CHAINID in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID.log"
$BINARY start --log_level debug --home $CHAIN_DIR/$CHAINID --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT" --grpc-web.address="0.0.0.0:$GRPCWEB"