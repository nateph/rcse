package cmd

import (
	"errors"
	"fmt"
	"io"
	"rcse/cmd/cliconfig"
	"rcse/pkg/concurrent"
	"rcse/pkg/files"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	commandToRun string
)

// ShellOptions contains options for the shell command
type ShellOptions struct {
	CommandToRun       string
	FailureLimit       int
	Forks              int
	IgnoreHostKeyCheck bool
	InventoryFile      string
	ListHosts          bool
	Password           string
	User               string
}

func newShellCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "shell",
		Short:        "Execute a shell shell command",
		SilenceUsage: true,
		RunE:         runShell,
	}

	flags := cmd.Flags()
	cliSettings.AddFlags(flags)
	flags.StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")

	return cmd
}

func runShell(cmd *cobra.Command, args []string) error {
	if cliSettings.InventoryFile == "" {
		return errors.New("no inventory flag was specified, all rcse operations require an inventory")
	}
	if commandToRun == "" {
		return errors.New("no command was found to run. exiting")
	}
	// if --username and --password were supplied correctly without --list-hosts
	if cliSettings.User != "" && cliSettings.Password == "default" && !cliSettings.ListHosts {
		cliSettings.Password = cliconfig.CheckAndConsumePassword(cliSettings.User, cliSettings.Password)
	}

	shellOptions := ShellOptions{
		CommandToRun:       commandToRun,
		FailureLimit:       cliSettings.FailureLimit,
		Forks:              cliSettings.Forks,
		IgnoreHostKeyCheck: cliSettings.IgnoreHostKeyCheck,
		InventoryFile:      cliSettings.InventoryFile,
		ListHosts:          cliSettings.ListHosts,
		Password:           cliSettings.Password,
		User:               cliSettings.User,
	}

	return executeShell(shellOptions)
}

func executeShell(shellOptions ShellOptions) error {
	parsedInventoryFile, err := files.LoadInventory(shellOptions.InventoryFile)
	if err != nil {
		return err
	}

	parsedHosts := parsedInventoryFile.Hosts

	if shellOptions.ListHosts {
		for _, host := range parsedHosts {
			fmt.Println(host)
		}
		return nil
	}
	//---------------------------------
	jobs := createJobs(parsedHosts, shellOptions)

	p := concurrent.NewPool(jobs, shellOptions.Forks)
	p.Run()

	var numErrors int
	for _, task := range p.Jobs {
		if task.Err != nil {
			logrus.Error(task.Err)
			numErrors++
		}
		if numErrors >= shellOptions.FailureLimit {
			logrus.Errorf("Failure limit of %d reached. Stopping.", shellOptions.FailureLimit)
			break
		}
	}
	return nil
}

// createJobs will gather jobs by supplying a ShellOptions for each host
func createJobs(parsedHosts []string, shellOptions ShellOptions) []*concurrent.Job {
	var jobs []*concurrent.Job

	for _, host := range parsedHosts {
		shellCmdOpts := cliconfig.CommandOptions{
			Host:               host,
			CommandToRun:       commandToRun,
			Sudo:               false,
			IgnoreHostkeyCheck: shellOptions.IgnoreHostKeyCheck,
			User:               shellOptions.User,
			Password:           shellOptions.Password,
		}
		jobs = append(jobs, concurrent.NewJob(cliconfig.RunCommand, shellCmdOpts))
	}

	return jobs
}

// 	// ---------------------------------------------
// 	results := make(chan cliconfig.CommandResult, shellOptions.Forks)
// 	timeout := time.After(15 * time.Second)
// 	jobs := make(chan cliconfig.CommandOptions, shellOptions.Forks)
// 	errorsChan := make(chan error)

// 	// Spawn x number of workers specified by --forks
// 	for w := 0; w < shellOptions.Forks; w++ {
// 		go worker(jobs, results)
// 	}

// 	go func() {
// 		for _, host := range parsedHosts {
// 			shellCmdOpts := cliconfig.CommandOptions{
// 				Host:               host,
// 				CommandToRun:       commandToRun,
// 				Sudo:               false,
// 				IgnoreHostkeyCheck: shellOptions.IgnoreHostKeyCheck,
// 				User:               shellOptions.User,
// 				Password:           shellOptions.Password,
// 			}
// 			jobs <- shellCmdOpts
// 		}
// 		close(jobs)
// 	}()

// 	for i := 0; i < len(parsedHosts); i++ {
// 		select {
// 		case res := <-results:
// 			res.PrintHostOutput()
// 		case <-timeout:
// 			fmt.Println("timed out")
// 		}
// 	}
// 	return nil
// }

// func worker(jobs <-chan cliconfig.CommandOptions, results chan<- cliconfig.CommandResult, errorsChan chan<- error) {
// 	for job := range jobs {
// 		fmt.Printf("job %v started\n", job)
// 		results <- job.RunCommands()
// 		fmt.Printf("job %v finished\n", job)
// 	}
// }
