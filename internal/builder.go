// Package internal provides core functionality for building and managing
// triggers and jobs, including runtime signal handling and configuration management.
package internal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Dr-Deep/hbsdsrv-build/internal/config"
	"github.com/Dr-Deep/logging-go"
)

// Builder manages triggers and jobs, handles runtime signals,
// and coordinates execution flow with logging and configuration support.
type Builder struct {
	trigger []Trigger
	jobs    map[string][]Job

	Logger *logging.Logger
	cfg    *config.Configuration

	currentRunningJob Job

	// Runtime Signals
	triggersSignalChan chan TriggerSignal
	exitSignalChan     chan os.Signal

	sync.Mutex
}

// NewBuilder creates and initializes a new Builder instance
// with the provided logger and configuration.
func NewBuilder(logger *logging.Logger, cfg *config.Configuration) *Builder {
	var builder = &Builder{
		trigger: []Trigger{},
		jobs:    map[string][]Job{},
		Logger:  logger,
		cfg:     cfg,
	}

	builder.exitSignalChan = make(chan os.Signal, 1)
	builder.triggersSignalChan = make(chan TriggerSignal, 1)

	return builder
}

// RegisterTrigger registers a new trigger with the Builder,
// using the provided trigger creation function and configuration.
func (b *Builder) RegisterTrigger(f func(config.TriggerConfig) Trigger, t config.TriggerConfig) {
	b.trigger = append(
		b.trigger,
		f(t),
	)
	b.Logger.Debug(
		fmt.Sprintf("%v", t),
	)
}

// RegisterJob registers a job under the given job name
// and associates it with a specific target.
func (b *Builder) RegisterJob(jobname string, targetname string, target Job) {
	b.jobs[jobname] = append(
		b.jobs[jobname],
		target,
	)

	b.Logger.Debug(
		fmt.Sprintf("%s, %s, %#v", jobname, targetname, target),
	)
}

// Launch starts the Builder by listening for system signals
// and executing triggers in separate goroutines.
func (b *Builder) Launch() {
	signal.Notify(
		b.exitSignalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	b.run()
}

// Stop gracefully stops the Builder by aborting the current job,
// stopping signal listeners, and closing channels before exiting.
func (b *Builder) Stop() {
	b.Logger.Info("Quitting...")

	// abort job
	if b.currentRunningJob != nil {
		b.currentRunningJob.Abort(b)
	}

	// wait for jobs
	b.Lock()

	// Stop Jobs
	signal.Stop(b.exitSignalChan)
	close(b.exitSignalChan)
	close(b.triggersSignalChan)

	os.Exit(0)
}

// RunJob executes the provided job within the Builder context,
// logging its execution and handling any errors.
func (b *Builder) RunJob(t *TriggerSignal, job Job) {
	b.Lock()
	b.Logger.Info(
		"Running Job",
		t.JobName, t.Reason,
	)
	//? ---

	b.currentRunningJob = job

	// run
	var err = job.Run(b)
	if err != nil {
		b.Logger.Error(
			t.JobName,
			err.Error(),
		)
	}

	b.currentRunningJob = nil

	//? ---
	b.Unlock()
	b.Logger.Info(
		"Completed Job",
		t.JobName, t.Reason,
	)
}

func (b *Builder) AbortCurrentJob() {
	if b.currentRunningJob == nil {
		return
	}

	if err := b.currentRunningJob.Abort(b); err != nil {
		//!log
		return
	}

	//!log
}

// ? fifo queue?
func (b *Builder) handleTrigger(t *TriggerSignal) {
	jobs := b.jobs[t.JobName]
	for _, j := range jobs {
		go b.RunJob(t, j)
	}
}

// queue, weil immer nur 1x job aufeinmal
func (b *Builder) run() {
	defer b.Stop()

	// launch trigger in go routines
	for _, t := range b.trigger {
		go t.Run(b, b.triggersSignalChan)
	}

	b.Logger.Info("triggers launched")
	for {
		select {
		case t := <-b.triggersSignalChan:
			b.handleTrigger(&t)

		case <-b.exitSignalChan:
			b.Logger.Info("catched SIGINT/SIGTERM")
			// b.Stop() in defer
			return

		default:
			continue
		}
	}
}
