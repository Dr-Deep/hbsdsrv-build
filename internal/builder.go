// Package internal provides core functionality for building and managing
// triggers and jobs, including runtime signal handling and configuration management.
package internal

import (
	"fmt"
	"hbsdsrv-build/internal/config"
	"hbsdsrv-build/internal/logging"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Builder manages triggers and jobs, handles runtime signals,
// and coordinates execution flow with logging and configuration support.
type Builder struct {
	trigger []Trigger
	jobs    map[string][]Job

	Logger *logging.Logger
	cfg    *config.Configuration

	// Runtime Signals
	triggersSignalChan chan TriggerSignal
	exitSignalChan     chan os.Signal

	sync.Mutex
}

// create a new Builder{} struct
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

func (b *Builder) RegisterTrigger(f func(config.TriggerConfig) Trigger, t config.TriggerConfig) {
	b.Logger.Debug(
		"registering",
		fmt.Sprintf("%#v", t),
	)
	b.trigger = append(
		b.trigger,
		f(t),
	)

	b.Logger.Debug(
		"registered",
		fmt.Sprintf("%#v", t),
	)
}

func (b *Builder) RegisterJob(jobname string, targetname string, target Job) {
	// []trigger
	// trigger-job: pkgbase/ports

	// JobPorts(targets)
	// targets = b.cfg.Jobs["trigger-job"]

	// trigger-job:
	// map[string]Targets
	// map[string][]Job

	//builder: jobs    map[string][]Job

	b.jobs[jobname] = append(
		b.jobs[jobname],
		target,
	)

	b.Logger.Debug(
		fmt.Sprintf("%s, %s, %#v", jobname, targetname, target),
	)
}

func (b *Builder) Launch() {
	signal.Notify(
		b.exitSignalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	b.run()
}

func (b *Builder) Stop() {
	b.Logger.Info("Quitting...")

	// wait for jobs
	b.Lock()

	// Stop Jobs
	signal.Stop(b.exitSignalChan)
	close(b.exitSignalChan)
	close(b.triggersSignalChan)

	os.Exit(0)
}

func (b *Builder) RunJob(job Job) {
	b.Lock()

	b.Logger.Info(
		"Got triggered, running job",
		fmt.Sprintf("%v", job),
	)

	if err := job.Run(b); err != nil {
		b.Logger.Error(
			fmt.Sprintf("%s", err),
			fmt.Sprintf("%#v", job),
		)
	}

	b.Unlock()
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

			return

		default:
			continue
		}
	}
}

func (b *Builder) handleTrigger(t *TriggerSignal) {
	jobs := b.jobs[t.JobName]
	for _, j := range jobs {
		go b.RunJob(j)
	}
}
