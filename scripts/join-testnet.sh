#!/usr/bin/env bash

set -euo pipefail

# This script will join the testnet and start the node

declare GOA_VERSION="v0.0.1-goa"
declare readonly BIN_PATH=""
declare readonly GITHUB_REPO="terra-money/alliance"
declare readonly GITHUB_URL="https://github.com/${GITHUB_REPO}"
declare readonly GITHUB_RAW="https://raw.githubusercontent.com/${GITHUB_REPO}/${GOA_VERSION}"

download_binaries (){
    mkdir -p "${HOME}/bin"
    GOA_GZ="testnets-$(uname -s)-$(uname -m).tar.gz"
    GOA_DOWNLOAD="${GITHUB_URL}/releases/download/${GOA_VERSION}/${GOA_GZ}"
    echo "Downloading ${GOA_DOWNLOAD}"
    curl -sSL "${GOA_DOWNLOAD}" | tar -xz -C "${HOME}/bin"
}

verify_binary(){
    local binary=$1
    if [ ! -f "$binary" ]; then
        echo "Binary $binary does not exist"
        exit 1
    fi
    echo $binary
}

verify_chain_id (){
    local chain_id=$1
    case $chain_id in
        "atreides-1")
            echo $chain_id
        ;;
        "corrino-1")
            echo $chain_id
        ;;
        "harkonnen-1")
            echo $chain_id
        ;;
        "ordos-1")
            echo $chain_id
        ;;
        *)
            echo "Chain ID $chain_id is not supported"
            exit 1
        ;;
    esac
}

get_binary_name(){
    local chain_id=$1
    echo "$(cut -d "-" -f1 <<< $chain_id)d"
}

get_binary_path(){
    local chain_id=$1
    echo "$HOME/bin/$(get_binary_name $chain_id)"
}

get_prefix(){
    local chain_id=$1
    cut -d "-" -f1 <<< $chain_id
}

get_denom(){
    local chain_id=$1
    echo "u$(cut -c-3  <<< $chain_id)"
}

get_port_prefix(){
    local chain_id=$1
    case $chain_id in
        "atreides-1")
            echo 414
        ;;
        "corrino-1")
            echo 413
        ;;
        "harkonnen-1")
            echo 411
        ;;
        "ordos-1")
            echo 412
        ;;
    esac
}

get_peers(){
    for (( i=0; i<3; i++ )); do
        curl -sSL "https://${PREFIX}.terra.dev:26657/status" | \
        awk -vRS=',' -vFS='"' '/id":"/{print $4}; /listen_addr":"[0-9]/{print $4}' |\
        paste -sd "@" -
    done | paste -sd "," -
}

parse_options(){
    while [ $# -gt 0 ]; do
        case "$1" in
            -c|--chain-id)
                CHAIN_ID=$(verify_chain_id $2)
                PREFIX=$(get_prefix $CHAIN_ID)
                DENOM=$(get_denom $CHAIN_ID)
                BINARY=$(get_binary $CHAIN_ID)
                shift 2
                ;;
            -m|--moniker)
                MONIKER=$2
                shift 2
                ;;
            --)
                shift
                break
                ;;
            *)
                echo "Not implemented: $1" >&2
                exit 1
                ;;
        esac
    done
}

main(){
    parse_options $@
    rm -rf /Users/greg/.ordos

    echo "Downloading binaries... This may take some time"
    download_binaries
    
    echo "Initializing node"
    
    if [ ! -f "$HOM" ]; then 
        $HOME/bin/$BINARY init "${MONIKER}" --chain-id "${CHAIN_ID}" 2>&1 | sed -e 's/{.*}//' 
    fi

    echo "Downloading genesis file"
    curl -sSL "${GITHUB_RAW}/genesis/${CHAIN_ID}/genesis.json" -o "${HOME}/.${PREFIX}/config/genesis.json"

    echo "Getting peer list"
    PEERS="$(get_peers)"

    echo "Starting node"
    exec $HOME/bin/$BINARY start \
        --p2p.persistent_peers "$PEERS" 
}

main $@
