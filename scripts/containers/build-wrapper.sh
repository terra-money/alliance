#!/usr/bin/env sh
set -x

export PATH=$PATH:/allianced/build/allianced
BINARY=/allianced/build/allianced
ID=${ID:-0}
LOG=${LOG:-allianced.log}
LD_LIBRARY_PATH=/lib

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found."
	exit 1
fi

export ALLIANCEDHOME="/allianced/data/node${ID}/allianced"

if [ -d "$(dirname "${ALLIANCEDHOME}"/"${LOG}")" ]; then
  exec 1> >(tee "${ALLIANCEDHOME}/${LOG}")
fi

exec $@
