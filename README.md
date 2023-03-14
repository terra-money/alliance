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
    - [Join the testnet](https://alliance.terra.money/guides/join-the-testnet)
- Technical specifications:
    - [Module parameters](https://alliance.terra.money/tech/parameters) and [Alliance asset properties](https://alliance.terra.money/tech/asset)
    - [Txs and queries](https://alliance.terra.money/tech/tx-queries)
    - [Data structures](https://alliance.terra.money/tech/data)
    - [State transitions](https://alliance.terra.money/tech/transitions)
    - [Invariants](https://alliance.terra.money/tech/invariants)
    - [Benchmarks](https://alliance.terra.money/tech/benchmarks)

## The `x/alliance` module

The Alliance module can be [added to any compatible Cosmos chain](https://alliance.terra.money/guides/get-started) and does not require any changes to consensus or major changes to common core modules. This module wraps around a chain’s native staking module, allowing whitelisted assets to be staked and earn rewards. Alliance assets can be staked with the Alliance module, and the chain's native staking module is used for native stakers. 


Chains that want to add `x/alliance` must enable the following modules:

- [x/auth](https://github.com/cosmos/cosmos-sdk/blob/main/x/auth/README.md)
- [x/bank](https://github.com/cosmos/cosmos-sdk/blob/main/x/bank/README.md)
- [x/ibc](https://github.com/cosmos/ibc-go#ibc-go)
- [x/staking](https://github.com/cosmos/cosmos-sdk/blob/main/x/staking/README.md)
- [x/distribution](https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/README.md)
- [x/gov](https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/README.md)

## Development environment

This project uses [Go v1.19](https://go.dev/dl/) and was originally bootstrapped using [Ignite CLI v0.25.1](https://docs.ignite.com/). However, for ease of upgrade, ignite has been removed in favor of manual workflows.

To set up the local development environment, clone this repo and run the following command:

```
$ make serve
```

If you want to build a ready-to-use binary, run the following:

```
$ make install
```

To build the proto files:
```
$ make proto-gen
```

## Localnet 

You can use a Docker orchestration to create a local network with 3 Docker containers:

- **localnet-start**: stop the testnet if running, build the terra-money/localnet-alliance image and start the nodes.
- **localnet-alliance-rmi**: removes the previously created terra-money/localnet-alliance image.
- **localnet-build-env**: delete and rebuild the terra-money/localnet-alliance
- **localnet-build-nodes**: using the terra-money/localnet-alliance starts a 3 docker containers testnet.
- **localnet-stop**: stop the testnet if running.

## Join the testnet

Joining the testnet is a very standardized process cosmos chain. In this case you will have to use **allianced** and follow [Terra documentation](https://docs.terra.money/full-node/manage-a-terra-validator/) since it's the same process but replacing it's genesis with the one that you can find in this repo under the path [docs/testnet/genesis.json](docs/testnet/genesis.json) and the following [seeds](http://3.75.187.158:26657/net_info),


### Running the simulation

The simulation app does not run out of the box because the Alliance module owns all native stake. The `x/staking` module's operation.go file panics when a delegator does not have a private key.

In order to run the simulation, update the `x/staking` module directly before compiling the simulation app using the following command. 

```shell
go mod vendor
sed -i '' 's/fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)/nil/g' vendor/github.com/cosmos/cosmos-sdk/x/staking/simulation/operations.go
ignite chain simulate
```

## Warning

Please note that Alliance is still undergoing final testing before its official release. TFL does not give any warranties, whether express or implied, as to the suitability or usability of the software or any of its content.

TFL will not be liable for any loss, whether such loss is direct, indirect, special or consequential, suffered by any party as a result of their use of the software or content.

Should you encounter any bugs, glitches, lack of functionality or other problems on the website, please submit bugs and feature requests through Github Issues. 
