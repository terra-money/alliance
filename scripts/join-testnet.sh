#!/usr/bin/env bash

set -euo pipefail

# This script will join the testnet and start the node
declare readonly OS=$(uname -s)
declare readonly PLATFORM=$(uname -m)
declare readonly GO_VERSION="1.19.5"
declare readonly GO_MD5="09e7f3b3ef34eb6099fe7312ecc314be"
declare readonly GOA_VERSION="v0.0.1-goa"
declare readonly BIN_PATH="${HOME}/go/bin/"
declare readonly GITHUB_REPO="terra-money/alliance"
declare readonly GITHUB_URL="https://github.com/${GITHUB_REPO}"
declare readonly GITHUB_RAW="https://raw.githubusercontent.com/${GITHUB_REPO}/${GOA_VERSION}"

declare GTMPDIR="${TMPDIR:-/tmp}"

error(){
    echo "Error: $1"
    exit 1
}

install_prereqs(){
    if [ $OS == "Linux" ] && [ -n "$(which apt)" ]; then 
        sudo apt update -y
        sudo apt install -y build-essential
    elif [ $OS == "Linux" ] && [ -n "$(which yum)" ]; then
        sudo yum update -y
        sudo yum group install -y "Development Tools"
    else
        echo "WARNING: You may need to install the gcc compiler"
    fi
}

download_go (){
    local tmpdir=$1
    if [ $OS == "Linux" ] && [ $PLATFORM == "x86_64" ]; then
       GO_GZ="go${GO_VERSION}.linux-amd64.tar.gz" 
    elif [ $OS == "Darwin" ] && [ $PLATFORM == "arm64" ]; then
       GO_GZ="go${GO_VERSION}.darwin-arm64.tar.gz"
    else
        error "Unsupported OS/Platform"
    fi
    GO_DOWNLOAD="https://go.dev/dl/${GO_GZ}"
    cd ${GTMPDIR}
    if [ ! -f "${GO_GZ}" ]; then
        echo "Downloading go from ${GO_DOWNLOAD}"
        curl -L "${GO_DOWNLOAD}"  -o ${GO_GZ}
    fi
    echo "Extracting ${GO_GZ}"
    # need to check md5sum
    tar -xzf ${GO_GZ} -C "${tmpdir}"
    echo
}

download_source (){
    local tmpdir=$1
    GOA_GZ="${GOA_VERSION}.tar.gz" 
    GOA_DOWNLOAD="${GITHUB_URL}/archive/refs/tags/${GOA_GZ}"
    cd ${GTMPDIR}
    if [ ! -f "${GOA_GZ}" ]; then
        echo "Downloading game of alliance from ${GOA_DOWNLOAD}"
        curl -sSL "${GOA_DOWNLOAD}" -o ${GOA_GZ}
    fi
    # need to check md5sum
    echo "Extracting ${GOA_GZ}"
    tar -xzf ${GOA_GZ} -C "${tmpdir}"
    echo
}

create_binary(){
    local prefix=$1
    local binary="${prefix}d"
    local tmpdir=$(mktemp -d)
    download_go ${tmpdir}
    download_source ${tmpdir}
    cd ${tmpdir}/alliance*
    export PATH="${tmpdir}/go/bin:${PATH}"
    export GOROOT="${tmpdir}/go"
    echo "Building ${binary}..."
    mkdir -p "${BIN_PATH}"
    go build -mod=readonly \
        -tags "netgo,ledger" \
        -ldflags " \
        -X github.com/cosmos/cosmos-sdk/version.Name=${prefix} 
        -X github.com/cosmos/cosmos-sdk/version.AppName=${prefix} 
        -X github.com/terra-money/alliance/app.Bech32Prefix=${prefix}
        -X github.com/terra-money/alliance/app.AccountAddressPrefix=${prefix} 
        -X github.com/terra-money/alliance/app.Name=${prefix}
        -X 'github.com/cosmos/cosmos-sdk/version.BuildTags=netgo,ledger'" \
        -trimpath -o ${BIN_PATH}${binary} ./cmd/allianced
    echo "Binary is located at ${BIN_PATH}${binary}"
    rm -rf ${tmpdir}
}

verify_chain_id (){
    local chain_id=$1
    case $chain_id in
        "atreides-1")
            echo "Chain ID is set to $chain_id"
        ;;
        "corrino-1")
            echo "Chain ID is set to $chain_id"
        ;;
        "harkonnen-1")
            echo "Chain ID is set to $chain_id"
        ;;
        "ordos-1")
            echo "Chain ID is set to $chain_id"
        ;;
        "")
            error "Chain ID is not set"
        ;;
        *)
            error "Chain ID $chain_id is not supported"
        ;;
    esac
}

get_binary(){
    local chain_id=$1
    echo "$(cut -d "-" -f1 <<< $chain_id)d"
}


get_prefix(){
    local chain_id=$1
    cut -d "-" -f1 <<< $chain_id
}

get_denom(){
    local prefix=$1
    echo "u$(cut -c-3  <<< $prefix)"
}

get_moniker(){
    local prefix=$1
    local cfgdir="${HOME}/.${prefix}/config"
    local moniker_txt="${cfgdir}/moniker.txt"
    if [ ! -f "${moniker_txt}" ]; then
        mkdir -p "${cfgdir}"
        echo "${prefix}-$(uuidgen | tr '[:upper:]' '[:lower:]' | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)" > "${moniker_txt}"
    fi
    cat "${moniker_txt}"
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
    set +u
    while [ $# -gt 0 ]; do
        case "$1" in
            -b|--binary)
                BINARY=$2
                shift 2
                ;;
            -c|--chain-id)
                CHAIN_ID=$2
                shift 2
                ;;
            -d|--denom)
                DENOM=$2
                shift 2
                ;;
            -m|--moniker)
                MONIKER=$2
                shift 2
                ;;
            -p|--prefix)
                PREFIX=$2
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
    verify_chain_id "${CHAIN_ID}"
    if [ -z "${BINARY}" ]; then
        BINARY=$(get_binary $CHAIN_ID)
    fi
    if [ -z "${PREFIX}" ]; then
        PREFIX=$(get_prefix $CHAIN_ID)
    fi
    if [ -z "${DENOM}" ]; then
        DENOM=$(get_denom $PREFIX)
    fi
    if [ -z "${MONIKER}" ]; then
        MONIKER=$(get_moniker $PREFIX)
    fi
    set -u
}

docs_bin_path(){
    echo "####################################################################"
    echo "# Binary is located in:"
    echo "####################################################################"
    echo
    echo "${BIN_PATH}${BINARY}"
    echo
}

docs_start_chain(){
    echo "####################################################################"
    echo "# To start the chain, run the following command:"
    echo "####################################################################"
    echo
    echo "${BIN_PATH}${BINARY} start --p2p.persistent_peers ${PEERS}" 
    echo
}

docs_create_wallet(){
    echo "####################################################################"
    echo "# To create a wallet run the following command:"
    echo "####################################################################"
    echo
    echo "${BIN_PATH}${BINARY} keys add ${MONIKER}" 
    echo
}

docs_faucet(){
    echo "####################################################################"
    echo "# To get tokens from the faucet visit:"
    echo "####################################################################"
    echo 
    echo "visit: https://game-of-alliances.terra.money/faucet"
    echo
}

docs_validate(){
    echo "####################################################################"
    echo "# To create a validator run the following command:"
    echo "####################################################################"
    echo
    echo "${BIN_PATH}${BINARY} tx staking create-validator \\"
    echo "  --amount=10000000${DENOM} \\"
    echo "  --pubkey='$(${BINARY} tendermint show-validator)' \\"
    echo "  --moniker="${MONIKER}" \\"
    echo "  --chain-id="${CHAIN_ID}" \\"
    echo "  --commission-rate="0.10" \\"
    echo "  --commission-max-rate="0.20" \\"
    echo "  --commission-max-change-rate=\"0.01\" \\"
    echo "  --min-self-delegation="1" \\"
    echo "  --gas=\"auto\" \\"
    echo "  --gas-adjustment=\"1.5\" \\"
    echo "  --gas-prices=\"0.025${DENOM}\" \\"
    echo "  --from=\"${MONIKER}\" \\"
    echo "  --yes"
    echo
}

main(){
    parse_options $@
    install_prereqs
    
    # prepend local path
    PATH="${BIN_PATH}:${PATH}"
    
    # check if binary exists
    if [ -z "$(which ${BINARY})" ]; then 
        create_binary "${PREFIX}"
    fi
    
    echo "Initializing node"
    
    if [ ! -f "${HOME}/.${PREFIX}/config/genesis.json" ]; then
        $BINARY init "${MONIKER}" --chain-id "${CHAIN_ID}" 2>&1 | sed -e 's/{.*}//' 
    fi

    echo "Downloading genesis file"
    curl -sSL "${GITHUB_RAW}/genesis/${CHAIN_ID}/genesis.json" -o "${HOME}/.${PREFIX}/config/genesis.json"

    echo "Getting peer list"
    PEERS="$(get_peers)"

    echo
    docs_bin_path
    docs_create_wallet
    docs_faucet
    docs_validate
    docs_start_chain
}

main $@
