package internal

import "fmt"

type Job interface {
	Run(b *Builder) error
	Abort(b *Builder) error
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
	err := job.Run(b)
	b.currentRunningJob = nil

	//? ---
	b.Unlock()
	if err != nil {
		b.Logger.Error(
			"Completed Job with Error",
			t.JobName, t.Reason,
			err.Error(),
		)
	} else {
		b.Logger.Info(
			"Completed Job",
			t.JobName, t.Reason,
		)
	}
}

func (b *Builder) AbortCurrentJob() {
	if b.currentRunningJob == nil {
		return
	}

	if err := b.currentRunningJob.Abort(b); err != nil {
		b.Logger.Error(err.Error())
		return
	}

	b.Logger.Info("Aborted current Job (clean)")
}
