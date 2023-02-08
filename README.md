<p align="center">
<h1 align="center"> 🤝 Alliance</h1>

<p align="center">
  Litepaper
  ·
  <a href="https://alliance.terra.money/">Technical Documentation</a>
  ·
  <a href="https://alliance.terra.money/guides/get-started">Integration Guide</a>
</p>

<br/>

# x/alliance interchain security

The Alliance module is part of the Interchain Security (Cosmos Shared Security that benefits from the IBC standard). Alliance is a friction free Interchain Security solution because there is no necessity to share hardware resources, have the blockchains synchronized nor modify the core of the origin chain that provide Interchain Security. Alliance module introduces the concept of alliance coins that can be seen as foreign coins bridged thru an IBC channel (ICS-004), whitelisted with the help of on-chain governance in the Alliance module and delegated by users or smart contracts to the active set of network validators.

Delegators of the alliance coins will be subjected to similar rules as the delegators of native coins but these delegators will provide Interchain Security to the network. The previously mentioned foreign coins can be in the form of Liquid Staked Derivative that benefits from the inflation of its native chain or any type of coin that can be bridged through the previously mentioned channels respecting the ICS-020 standard.

When users delegate coins through the Alliance module the voting power of the validators will be diluted, as a consequence, bad actors will have to increase their capital spendings to try and corrupt the consensus of the blockchain. 

By design, x/alliance use the following CosmosSDK modules to implement interchain security to a new or existing blockchain :

- [x/auth](https://github.com/cosmos/cosmos-sdk/blob/main/x/auth/README.md),
- [x/bank](https://github.com/cosmos/cosmos-sdk/blob/main/x/bank/README.md),
- [x/ibc](https://github.com/cosmos/ibc-go#ibc-go),
- [x/staking](https://github.com/cosmos/cosmos-sdk/blob/main/x/staking/README.md), 
- [x/distribution](https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/README.md), 
- [x/gov](https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/README.md).


# Development environment
This project uses [Go v1.19](https://go.dev/dl/) and was bootstrapped with [Ignite CLI v0.25.1](https://docs.ignite.com/). However, for ease of upgrade, ignite has been removed in favor of manual workflows.

To run the local development environment use:
```
$ make serve
```

If you want to build a binary ready to use:
```
$ make install
```

To build the proto files:
```
$ make proto-gen
```

## Localnet 
Docker orchestration to create a local network with 3 docker containers:

- **localnet-start**: stop the testnet if running, build the terra-money/localnet-alliance image and start the nodes.
- **localnet-alliance-rmi**: removes the previously created terra-money/localnet-alliance image.
- **localnet-build-env**: delete and rebuild the terra-money/localnet-alliance
- **localnet-build-nodes**: using the terra-money/localnet-alliance starts a 3 docker containers testnet.
- **localnet-stop**: stop the testnet if running.

## Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

## Join the Game of Alliance Testnet
> Note: This is only for advanced users. It assumes that you already know how to run a full node.

In order to join the testnet, you will have to build your node to match any one of the running chains. 

1. Pick one chain: `ordos | corrino | harkonnen | atreides`
2. Run `make build-alliance ACC_PREFIX=ordos` to build the binary for ordos. Replace ordos with other prefixes to build for the other chains.
3. Run `build/ordosd init <moniker>` to init the chain
4. Copy the genesis file from the `genesis/ordos-1` folder for the respective chain into `~/.ordos/config/` folder
5. Update the seed from the table below in `~/.ordos/config/config.yaml`
6. Run the chain with `build/ordosd start`
7. Follow the instructions on [setting up a validator](https://docs.terra.money/full-node/manage-a-terra-validator/set-up-validator)

| Chain       | Seeds                                                        |
|-------------|--------------------------------------------------------------|
| atreides-1  | 6ebf0000ee85ff987f1d9de3223d605745736ca9@35.168.16.221:41356 |
| corrino-1   | TBC                                                          |
| harkonnen-1 | TBC                                                          |
| ordos-1     | TBC                                                          |

### Running the simulation
The simulation app does not run out of the box since the alliance module owns all native stake. The `x/staking` module's operation.go file panics when a delegator does not have a private key.

In order to run the simulation, you can update the `x/staking` module directly before compiling the simulation app using the following command
```shell
go mod vendor
sed -i '' 's/fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)/nil/g' vendor/github.com/cosmos/cosmos-sdk/x/staking/simulation/operations.go
ignite chain simulate
```

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)

## Warning

Please note that this is a beta version of the which still undergoing final testing before its official release. TFL does not give any warranties, whether express or implied, as to the suitability or usability of the software or any of its content.

TFL will not be liable for any loss, whether such loss is direct, indirect, special or consequential, suffered by any party as a result of their use of the software or content.

Should you encounter any bugs, glitches, lack of functionality or other problems on the website, please submit bugs and feature requests through Github Issues. 
