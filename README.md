<p align="center">
<h1 align="center"> 🤝 Alliance</h1>

<p align="center">
  <a href="https://alliance.terra.money/">Technical Documentation</a>
  ·
  <a href="https://alliance.terra.money/guides/get-started">Integration Guide</a>
  ·
    <a href="https://alliance.terra.money/alliance-audit.pdf">Code Audit</a>

</p>

<br/>

## Overview

Alliance is an open-source Cosmos SDK module that leverages interchain staking to form ***economic alliances*** among blockchains. By boosting the economic activity across Cosmos chains through creating bilateral, mutually beneficial alliances, Alliance aims to give rise to a new wave of innovation, user adoption, and cross-chain collaboration.

**Alliance allows blockchains to trade yield with each other- think of it like yield farming for L1s.**

### Here’s how it works:

- Two chains [integrate the Alliance module](https://alliance.terra.money/guides/get-started) and decide through governance which assets can be staked on their chain. These are known as [Alliance assets](https://alliance.terra.money/alliance#what-are-alliance-assets). 
- Each Alliance asset is assigned a [Take Rate](https://alliance.terra.money/alliance#the-take-rate) (the percentage of staked Alliance assets the chain redistributes to native chain stakers) and a [Reward Weight](https://alliance.terra.money/alliance#rewards) (the percentage of native staking rewards the chain distributes to Alliance asset stakers).
- Users of each chain can then [bridge their assets via IBC](https://alliance.terra.money/alliance#what-are-alliance-assets) to the other chain and stake them to earn the Reward Weight. 

## Tech Specs

[The Alliance Docs](https://alliance.terra.money/) contain detailed information about Alliance. Familiarize yourself with the following concepts before integrating the Alliance module. 

- About Alliance
    - [Overview](https://alliance.terra.money/overview)
    - [How Alliance works](https://alliance.terra.money/alliance)
    - [Alliance staking](https://alliance.terra.money/concepts/staking)
    - [Reward distribution](https://alliance.terra.money/concepts/rewards)
    - [Validator shares](https://alliance.terra.money/concepts/delegation)
- Guides
    - [Integrate the Alliance module](https://alliance.terra.money/guides/get-started)
    - [Create, update, or delete an Alliance](https://alliance.terra.money/guides/create)
    - [Interact with an Alliance](https://alliance.terra.money/guides/how-to)
- Technical specifications:
    - [Module parameters](https://alliance.terra.money/tech/parameters) and [Alliance asset properties](https://alliance.terra.money/tech/asset)
    - [Txs and queries](https://alliance.terra.money/tech/tx-queries)
    - [Data structures](https://alliance.terra.money/tech/data)
    - [State transitions](https://alliance.terra.money/tech/transitions)
    - [Invariants](https://alliance.terra.money/tech/invariants)
    - [Benchmarks](https://alliance.terra.money/tech/benchmarks)

## Integrate the `x/alliance` module

The Alliance module can be [added to any compatible Cosmos chain](https://alliance.terra.money/guides/get-started) and does not require any changes to consensus or major changes to common core modules. This module wraps around a chain’s native staking module, allowing whitelisted assets to be staked and earn rewards. Alliance assets can be staked with the Alliance module, and the chain's native staking module is used for native stakers. 

Chains that want to add `x/alliance` must enable the following modules:

- [x/auth](https://github.com/cosmos/cosmos-sdk/blob/main/x/auth/README.md)
- [x/bank](https://github.com/cosmos/cosmos-sdk/blob/main/x/bank/README.md)
- [x/ibc](https://github.com/cosmos/ibc-go#ibc-go)
- [x/staking](https://github.com/cosmos/cosmos-sdk/blob/main/x/staking/README.md)
- [x/distribution](https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/README.md)
- [x/gov](https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/README.md)

For an in-depth guide on integrating `x/alliance`, visit the [Alliance Module Integration Guide](https://alliance.terra.money/guides/get-started). 

## Development environment

The following sections are for developers working on the  `x/alliance` module. 

This project uses [Go v1.19](https://go.dev/dl/).

To build a ready-to-use binary, run the following:

```sh
make install
```

### Localnet

The Localnet is a development environment that uses a Docker orchestration to create a local network with 3 Docker containers:

- `make localnet-start` : stop the testnet if running, build the `terra-money/localnet-alliance` image and start the nodes.
- `make localnet-alliance-rmi`: remove the previously created `terra-money/localnet-alliance` image.
- `make localnet-build-env`: delete and rebuild the `terra-money/localnet-alliance`
- `make localnet-build-nodes`: using the `terra-money/localnet-alliance` starts a 3 docker containers testnet.
- `make localnet-stop`: stop the testnet if running.

### Running the simulation

The simulation app does not run out of the box because the Alliance module owns all the native stake. The `x/staking` module's `operation.go` file panics when a delegator does not have a private key.

Use the following command to update the `x/staking` module directly before compiling the simulation app.

```sh
go mod vendor
sed -i '' 's/fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)/nil/g' vendor/github.com/cosmos/cosmos-sdk/x/staking/simulation/operations.go
go test -benchmem -run=^$ -bench ^BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true
```

## Warning

Please note that Alliance is still undergoing final testing before its official release. TFL does not give any warranties, whether express or implied, as to the suitability or usability of the software or any of its content.

TFL will not be liable for any loss, whether such loss is direct, indirect, special or consequential, suffered by any party as a result of their use of the software or content.

Should you encounter any bugs, glitches, lack of functionality or other problems on the website, please submit bugs and feature requests through Github Issues.
