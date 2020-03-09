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
func Execute(hosts []string, opts cliconfig.JobOptions) error {
	jobs := make(chan command.Options, opts.Forks)
	results := make(chan command.Result, opts.Forks)
	errorsChan := make(chan error, len(hosts))

	var wg sync.WaitGroup
	// Spawn n number of workers specified by --forks
	for w := 0; w < opts.Forks; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, errorsChan)
	}

	go func() {
		for _, host := range hosts {
			jobOpts := command.Options{
				Host:               host,
				CommandToRun:       opts.CommandToRun,
				IgnoreHostkeyCheck: opts.IgnoreHostKeyCheck,
				User:               opts.User,
				Password:           opts.Password,
			}
			jobs <- jobOpts
		}
	}()

	go func() {
		wg.Wait()
		close(jobs)
	}()

	var failureLimit int

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			res.PrintHostOutput(opts.OutFormat)
		case <-errorsChan:
			failureLimit++
			if failureLimit >= opts.FailureLimit {
				return errors.New("too many failures, exiting")
			}
		}
	}

	return nil
}
