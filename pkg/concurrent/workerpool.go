package concurrent

import (
	"rcse/cmd/cliconfig"
	"sync"
)

// Pool is a worker group that runs a number of jobs at a
// configured number of forks.
type Pool struct {
	Jobs []*Job

	forks    int
	jobsChan chan *Job
	wg       sync.WaitGroup
}

// NewPool initializes a new pool with the given jobs and
// at the given fork amount.
func NewPool(jobs []*Job, forks int) *Pool {
	return &Pool{
		Jobs:     jobs,
		forks:    forks,
		jobsChan: make(chan *Job),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 0; i < p.forks; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Jobs))
	for _, job := range p.Jobs {
		p.jobsChan <- job
	}

	// all workers return
	close(p.jobsChan)

	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool) work() {
	for job := range p.jobsChan {
		job.Run(&p.wg)
	}
}

// Job encapsulates a work item that should go in a work
// pool.
type Job struct {
	// Err holds an error that occurred during a job. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Err error

	f    func(cliconfig.CommandOptions) error
	args cliconfig.CommandOptions
}

// NewJob initializes a new job based on a given work
// function.
func NewJob(f func(cliconfig.CommandOptions) error, args cliconfig.CommandOptions) *Job {
	return &Job{f: f, args: args}
}

// Run runs a Job and does appropriate accounting via a
// given sync.WorkGroup.
func (j *Job) Run(wg *sync.WaitGroup) {
	// fmt.Printf("starting job %v\n", j)
	j.Err = j.f(j.args)
	// fmt.Printf("finished job %v\n", j)
	wg.Done()
}
