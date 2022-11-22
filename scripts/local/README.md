# Scripts

Folder containing some scripts to test or/and demo the functionality. Before executing any script you must have the `allianced` installed.

1. **init.sh** create the chain with the initial values
2. **start.sh** start the chain
3. **gov.sh** submit the gov.json governance proposal, votes on favor and query the created alliance
4. **delegate.sh** delegate to the previously create alliance and query the modified alliance
5. **rewards.sh** claim rewards and query information about the evidences of the process
6. **undelegate.sh** undelegante the tokens from the alliance and query the evidences
7. **gov-del.sh** with the file gov-delete.json deletes the alliance created in third step.

> This scripts must be executed in the specified order since they have dependencies on each other.
