[![Go Report Card](https://goreportcard.com/badge/github.com/nateph/rcse)](https://goreportcard.com/report/github.com/nateph/rcse ) ![GitHub](https://img.shields.io/github/license/nateph/rcse)
### RCSE (Run Command Somewhere Else)
***
`rcse` is meant to be a simple tool for remote machine automation written in Go, and using SSH under the hood, with no agent to install on any remote machines.

The idea is to have subcommands, i.e. `raw`, `shell`, `yum`, `playbook`, dictate what the program will do. The inventory file, supplied by `-i` will supply the list of hosts that the program is ran on, and it needs to be in yaml format. 

Example: 
```
hosts:
  - myhost001.ci.com
  - myhost002.ci.com
  - vmhost001.ci.com
```

For example, running `rcse raw -i <inventory_file> -c "ls -l"` would list the contents of your user on each of the hosts in the inventory file. Without a user specified, rcse will use the current user's ssh keys. 

A username, `--user`, and password, `--password` can be specified to execute the commands as a different user.
If using a different user, the program will prompt for password input, or you may pass the password through command line argument, but that is not recommended. 

`rcse ... -u root -p` and `rcse ... -u root -p=CoolPassword` are equivalent, but the first one prompts you for password input.

Each subcommand has its own set of flags relevant to its purpose, as well as global flags.

### Installation 
***
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