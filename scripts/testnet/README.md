# Scripts

This folder contains a sequence of helper scripts for creating an alliance on testnet and automatic delegation. 

1. **gov.sh** submits a gov.json governance proposal, votes in favor of it and then queries the created alliance.
2. **delegate.sh** delegates to the previously create alliance and queries the modified alliance.
3. **rewards.sh** claims available rewards and retrieves information about the process

> Note that these scripts must be executed in the specified order since they have dependencies on each other.
