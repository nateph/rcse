[![CircleCI](https://circleci.com/gh/circleci/circleci-docs.svg?style=shield)](https://circleci.com/gh/nateph/rcse) [![Go Report Card](https://goreportcard.com/badge/github.com/nateph/rcse)](https://goreportcard.com/report/github.com/nateph/rcse ) ![GitHub](https://img.shields.io/github/license/nateph/rcse)
### RCSE (Run Command Somewhere Else)
`rcse` is meant to be a simple tool for remote machine automation written in Go, and using SSH under the hood, with no agent to install on any remote machines.

The program has subcommands, i.e. `shell`, `yum`, `sequence`, and those will dictate what actions get performed.
#### Inventory
The inventory file, passed with `-i`, will supply the list of hosts that the program is ran on, and it needs to be in yaml format under the key `hosts`:
 
```
hosts:
  - myhost001.ci.com
  - myhost002.ci.com
  - vmhost001.ci.com
```

For example, running `rcse shell -i <inventory_file> -c "ls -l"` would list the contents of your user's home directory on each of the hosts in the inventory file. The command should be quoted, and without a user specified, it will use the current user's ssh keys. 

#### Authentication 
A username, `--user`, and password, `--password`, can be specified to execute the commands as a different user.
If providing a different user, the program will prompt for password input. You may pass the password explicitly in the command, however that is not recommended.

`rcse ... -u root -p` and `rcse ... -u root -p=CoolPassword` are equivalent, but the first one prompts you for password input securely. Note that passing a user without a password flag will still prompt for a password.

If no user is specified, it will use the current user's id_rsa found in `~/.ssh/id_rsa`.

#### Running concurrently 
There are options for how many hosts to run on at one time, with the `--forks` flag. Accompanying it is `--failure-limit` which will stop and exit the program when that limit is hit. With neither flag set, the forks and failure limit are set to 1 and 1000 respectively, making a default run safe and unopinionated about the failure limit.

At any one time, the program will keep `--forks=n` amount of executors running. For example, if set to 15, the program will start by executing on 15 hosts concurrently, and then when one finishes, it will pick up the next host to run on, and so on. Contrast this with batching, where the next batch will only start when all hosts in the previous batch have finished. Setting this flag executes on hosts as they are listed in the inventory file, in descending order.

#### Running a Sequence
The subcommand `sequence` allows for multiple commands to run on each host, in order, on each host in the inventory. A sequence is contained in a yaml file and looks like so:
```
jobs:
  - name: check red hat version
    module: shell
    command: cat /etc/redhat-release
  - name: print id
    module: shell
    command: id
```
Running this would look like `rcse sequence -f <sequence_file> -i <inventory_file>`. All options available to normal commands are available for a sequence as well, such as running with a different user, or adding forks and failure limits. It is important to note that currently the failure limit counter is per command, not per host.

#### Flags
Each subcommand has its own set of flags relevant to its purpose, as well as global flags available to every command.
Please use `--help` on any command to see options.

### Installation 
#### Download
Coming soon 
#### Using Go Modules
```
git clone git@github.com:nateph/rcse.git
cd rcse && make install
```
#### Using $GOPATH
```
go get github.com/nateph/rcse
cd $GOPATH/src/github.com/nateph/rcse && make install
```
See `Makefile` for more options