# Example Scripts

The following scripts are prepared to interact with the alliance localnet. To start the localnet run `make localnet-start` inside the alliance root folder. 

> If you are using linux take in consideration to change the `.testnets` folder owner to your host machine user by using `sudo chown [HOST_USER_NAME]:[HOST_USER_NAME] -R .testnets/` because the folder may be owned by the root user.

1. Create Alliance: create an alliance with the default token `stake` and vote for the proposal to pass
2. Delegate: delegate to the first validator
3. Undelegate: undelegate from the first validator
4. Redelegate: redelegate from the first validator to the second validator