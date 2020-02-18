### RCSE (Run Command Somewhere Else)
rcse is meant to be a simple tool for remote machine automation written in Go, and using SSH under the hood, with no agent to install on any remote machines.

The idea is to have subcommands, i.e. "shell" dictate what the program will do. So running `rcse raw_shell -i <inventory_file> -c "ls -ltr"` would execute the specified shell command on each of the hosts in the inventory file.

Each subcommand has its own set of flags relevant to its purpose, as well as global flags.
