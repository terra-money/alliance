# x/alliance interchain security

Independent blockchain module compatible with [CosmosSDK v0.46.x](https://github.com/cosmos/cosmos-sdk) that allow blockchains to share security without  necessity of modifying the core Cosmos SDK modules nor share any hardware resources. By design, x/alliance use the following CosmosSDK modules to implement interchain security to a new or existing blockchain :

- [x/auth](https://github.com/cosmos/cosmos-sdk/blob/main/x/auth/README.md),
- [x/bank](https://github.com/cosmos/cosmos-sdk/blob/main/x/bank/README.md),
- [x/ibc](https://github.com/cosmos/ibc-go#ibc-go),
- [x/staking](https://github.com/cosmos/cosmos-sdk/blob/main/x/staking/README.md), 
- [x/distribution](https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/README.md), 
- [x/gov](https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/README.md).

Since security on Delegated Proof of Stake Chains is directly related to the [voting power of the validators](https://docs.tendermint.com/v0.34/tendermint-core/validators.html#validators) x/alliance enable staking of native tokens to increase the chain security and users that stake will also earn staking rewards.

# Development environment

This project uses [Go v1.18](https://go.dev/dl/) and was bootstrapped with [Ignite CLI v0.24.0](https://docs.ignite.com/). 

To run the local development environment use:
```
$ ignite chain serve --verbose
```

If you want to build a binary ready to use:
```
$ ignite chain build
```

To build the proto files:
```
$ ignite generate proto-go
```


# Local Testnet 

Docker orchestration to create a local testnet with 4 different docker images:

- **localnet-start**: stop the testnet if running, build the terra-money/localnet-alliance image and start the nodes.
- **localnet-alliance-rmi**: removes the previously created terra-money/localnet-alliance image.
- **localnet-build-env**: delete and rebuild the terra-money/localnet-alliance
- **localnet-build-nodes**: using the terra-money/localnet-alliance starts a 4 docker containers testnet.
- **localnet-stop**: stop the testnet if running.