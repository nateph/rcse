[![CircleCI](https://circleci.com/gh/circleci/circleci-docs.svg?style=shield)](https://circleci.com/gh/nateph/rcse) [![Go Report Card](https://goreportcard.com/badge/github.com/nateph/rcse)](https://goreportcard.com/report/github.com/nateph/rcse ) ![GitHub](https://img.shields.io/github/license/nateph/rcse)
### RCSE (Run Command Somewhere Else)
`rcse` is meant to be a simple tool for remote machine automation written in Go, utilizing only SSH, with no agent or extra packages to install on any remote machines.

The program has subcommands, e.g. `shell` and `script`, and those will dictate what actions get performed.
#### Running a single command
Running `rcse shell -c "ls -la" -i <inventory_file>` would run "ls -la", listing the contents of your user's home directory on each of the hosts in the inventory file. The command should be quoted, and without a user specified, it will use the current local user's ssh keys (id_rsa by default, see --help for more options). 

#### Running a script
The subcommand `script` allows for a script, from your local machine, to be ran remotely on all hosts in the inventory.

Running this would look like `rcse script my_script.sh -i <inventory_file>`. All options available to normal commands are available for a script as well, such as running with a different user, or adding forks and failure limits.

#### Inventory
The inventory file, passed with `-i`, will supply the list of hosts that the program is ran on. It takes in a plaintext 
file like the one shown below.

```
$ cat my_inventory
myhost001.ci.com
myhost[002:010].ci.com
myhost012.ci.com
vmhost001.ci.com
baremetal[01:03].company.location.com
```

Note that if the hosts are in a contiguous range, you may specify the range with host[start:end]... and rcse will expand it internally like so:
```
baremetal[01:03].company.location.com
```
gets expanded to 
```
baremetal01.company.location.com
baremetal02.company.location.com
baremetal03.company.location.com
```

#### Authentication 
A username, `--user`, and password, `--password`, can be specified to execute the commands as a different user.
If providing a different user, the program will prompt for password input. You may pass the password explicitly in the command, however that is not recommended.

`rcse ... -u root -p` and `rcse ... -u root -p=CoolPassword` are equivalent, but the first one prompts you for password input securely. Note that passing a user without a password flag will still prompt for a password.

If no user is specified, it will use the current user's id_rsa found in `~/.ssh/id_rsa`.

#### Running concurrently 
There are options for how many hosts to run on at one time, with the `--forks` flag. Accompanying it is `--failure-limit` which will stop and exit the program when that limit is hit. With neither flag set, the forks and failure limit are set to 1 and 1000 respectively, making a default run safe and unopinionated about the failure limit.

At any one time, the program will keep `--forks=n` amount of executors running. For example, if set to 15, the program will start by executing on 15 hosts concurrently, and then when one finishes, it will pick up the next host to run on, and so on. Contrast this with batching, where the next batch will only start when all hosts in the previous batch have finished. Setting this flag executes on hosts as they are listed in the inventory file, in descending order.

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
