#!/usr/bin/env sh
set -euo pipefail
set -x

BINARY=/allianced/${BINARY:-allianced}
ID=${ID:-0}
LOG=${LOG:-allianced.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found."
	exit 1
fi

export ALLIANCEDHOME="/data/node${ID}/allianced"

if [ -d "$(dirname "${ALLIANCEDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${ALLIANCEDHOME}" "$@" | tee "${ALLIANCEDHOME}/${LOG}"
else
  "${BINARY}" --home "${ALLIANCEDHOME}" "$@"
fi
