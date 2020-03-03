[![CircleCI](https://circleci.com/gh/circleci/circleci-docs.svg?style=shield)](https://circleci.com/gh/nateph/rcse) [![Go Report Card](https://goreportcard.com/badge/github.com/nateph/rcse)](https://goreportcard.com/report/github.com/nateph/rcse ) ![GitHub](https://img.shields.io/github/license/nateph/rcse)
### RCSE (Run Command Somewhere Else)
`rcse` is meant to be a simple tool for remote machine automation written in Go, and using SSH under the hood, with no agent to install on any remote machines.

The program has subcommands, i.e. `raw`, `shell`, `yum`, `sequence`, and those will dictate what actions get performed.
#### Inventory
The inventory file, supplied by `-i` will supply the list of hosts that the program is ran on, and it needs to be in yaml format under the key `hosts`:
 
```
hosts:
  - myhost001.ci.com
  - myhost002.ci.com
  - vmhost001.ci.com
```

For example, running `rcse raw -i <inventory_file> -c "ls -l"` would list the contents of your user on each of the hosts in the inventory file. The command should be quoted, and without a user specified, it will use the current user's ssh keys. 

#### Authentication 
A username, `--user`, and password, `--password`, can be specified to execute the commands as a different user.
If providing a different user, the program will prompt for password input. You may pass the password explicitly in the command, however that is not recommended.

`rcse ... -u root -p` and `rcse ... -u root -p=CoolPassword` are equivalent, but the first one prompts you for password input securely. Note that passing a user without a password flag will still prompt for a password.

If no user is specified, it will use the current user's id_rsa found in `~/.ssh/id_rsa`.

Each subcommand has its own set of flags relevant to its purpose, as well as global flags available to every command.

### Installation 
#### Download
Coming soon 
#### Using Go Modules
```
git clone git@github.com:nateph/rcse.git
cd rcse && go install
```
#### Using $GOPATH
```
go get github.com/nateph/rcse
cd $GOPATH/src/github.com/nateph/rcse && go install
```