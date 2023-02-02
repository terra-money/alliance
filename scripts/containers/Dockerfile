# syntax=docker/dockerfile:1

ARG GO_VERSION="1.19"
ARG PLATFORM="bullseye"
ARG BUILDPLATFORM="linux/amd64"
ARG BASE_IMAGE="golang:${GO_VERSION}-${PLATFORM}"
ARG ACC_PREFIX=alliance

FROM --platform=${BUILDPLATFORM} ${BASE_IMAGE} as base

###############################################################################
# Builder
###############################################################################

FROM base as builder

ARG ACC_PREFIX
ARG BUILDPLATFORM

WORKDIR /go/src
COPY . .

RUN set -eux &&\
    apt update &&\
    apt install build-essential -y

RUN set -eux && sed \
    -e "s/AccountAddressPrefix = \"alliance\"/AccountAddressPrefix = \"${ACC_PREFIX}\"/g" \
    -e "s/Name                 = \"alliance\"/Name                 = \"${ACC_PREFIX}\"/g" \
    -i app/app.go

RUN set -eux &&\
    make build -o .

RUN set -eux &&\
    WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm | cut -d ' ' -f 2) && \
    cp /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@${WASMVM_VERSION}/internal/api/libwasmvm.$(uname -m).so /lib/

###### Node
FROM --platform=${BUILDPLATFORM} ${BASE_IMAGE} AS run

COPY --from=builder /go/bin/allianced /usr/local/bin/allianced
COPY --from=builder /lib/libwasmvm*.so /lib

WORKDIR /app
CMD ["allianced", "start", "--log_format", "plain"]
