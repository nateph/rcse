package concurrent

import (
	"errors"
	"sync"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/command"
	"github.com/sirupsen/logrus"
)

func worker(wg *sync.WaitGroup, jobs <-chan command.Options, results chan<- command.Result, errChan chan<- error) {
	defer wg.Done()
	for job := range jobs {
		res, err := job.RunCommand()
		if err != nil {
			logrus.Error(err)
			errChan <- err
			continue
		}
		results <- res
	}
}

// Execute is a wrapper function to handle concurrency
func Execute(conf *cliconfig.Config, inventory ...string) error {
	jobs := make(chan command.Options)
	results := make(chan command.Result)
	errorsChan := make(chan error)

	var wg sync.WaitGroup
	// Spawn n number of workers specified by --forks
	for w := 0; w < conf.Options.Forks; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, errorsChan)
	}

	generatedJobs := generateJobs(conf, inventory...)
	go func() {
		for _, j := range generatedJobs {
			jobs <- j
		}
	}()

	go func() {
		wg.Wait()
		close(jobs)
	}()

	var failureLimit int

	for i := 0; i < len(inventory)*len(conf.Jobs); i++ {
		select {
		case res := <-results:
			res.PrintHostOutput(conf.Options.OutFormat)
		case <-errorsChan:
			failureLimit++
			if failureLimit >= conf.Options.FailureLimit {
				return errors.New("too many failures, exiting")
			}
		}
	}

	return nil
}

// func generateShellJobs(s *cmd.ShellOptions, inventory ...string) []command.Options {
// 	var jobs []command.Options

// 	for _, host := range inventory {
// 		jobOpts := command.Options{
// 			Host:               host,
// 			Command:            s.Command,
// 			IgnoreHostkeyCheck: s.BaseOpts.IgnoreHostKeyCheck,
// 			User:               s.BaseOpts.User,
// 			Password:           s.BaseOpts.Password,
// 		}
// 		jobs = append(jobs, jobOpts)
// 	}
// 	return jobs
// }

func generateJobs(conf *cliconfig.Config, inventory ...string) []command.Options {
	var jobs []command.Options

	for _, host := range inventory {
		for _, job := range conf.Jobs {
			jobOpts := command.Options{
				Host:               host,
				Command:            job.Command,
				IgnoreHostkeyCheck: conf.Options.IgnoreHostKeyCheck,
				User:               conf.Options.User,
				Password:           conf.Options.Password,
			}
			jobs = append(jobs, jobOpts)
		}
	}
	return jobs
}
