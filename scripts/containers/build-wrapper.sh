#!/usr/bin/env sh
set -x

export PATH=$PATH:/allianced/build/allianced
BINARY=/allianced/build/allianced
ID=${ID:-0}
LOG=${LOG:-allianced.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found."
	exit 1
fi

export ALLIANCEDHOME="/allianced/data/node${ID}/allianced"

if [ -d "$(dirname "${ALLIANCEDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${ALLIANCEDHOME}" "$@" | tee "${ALLIANCEDHOME}/${LOG}"
else
  "${BINARY}" --home "${ALLIANCEDHOME}" "$@"
fi
