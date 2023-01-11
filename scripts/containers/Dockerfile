###### Builder
FROM golang:1.18.7-bullseye AS build

RUN apt update && apt install jq build-essential -y

RUN curl https://get.ignite.com/cli@v0.25.1! | bash

WORKDIR /allianced
COPY . /allianced

ARG ACC_PREFIX=alliance
ARG BUILDPLATFORM=linux/amd64

RUN sed -i "s/AccountAddressPrefix = \"alliance\"/AccountAddressPrefix = \"${ACC_PREFIX}\"/g" app/app.go
RUN sed -i "s/Name                 = \"alliance\"/Name                 = \"${ACC_PREFIX}\"/g" app/app.go

RUN ignite chain build --output . --verbose

RUN cp /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.`uname -m`.so /lib/

###### Node
FROM debian:bullseye AS run
COPY ./scripts/containers/build-wrapper.sh /allianced/build-wrapper.sh

VOLUME /allianced
COPY --from=build /allianced/ /allianced/
COPY --from=build /lib/ /lib/
WORKDIR /allianced

EXPOSE 26656 26657
ENTRYPOINT ["/allianced/build-wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
