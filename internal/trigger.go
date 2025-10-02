package internal

import (
	"fmt"

	"github.com/Dr-Deep/hbsdsrv-build/internal/config"
)

type Trigger interface {
	Run(b *Builder, c chan TriggerSignal)
}

type TriggerSignal struct {
	JobName string
	Reason  string
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

// ? fifo queue?
func (b *Builder) handleTrigger(t *TriggerSignal) {
	jobs := b.jobs[t.JobName]
	for _, j := range jobs {
		go b.RunJob(t, j)
	}
}
