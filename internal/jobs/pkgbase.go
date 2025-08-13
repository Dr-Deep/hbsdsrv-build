package jobs

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"os"
)

type JobPkgbase struct {
	args       config.JobArgs
	childProcs []*os.Process
}

func NewJobPkgbase(config.JobArgs) internal.Job {
	return &JobPkgbase{}
}

func (j *JobPkgbase) Run(b *internal.Builder) error {
	return nil
}

func (j *JobPkgbase) Abort(b *internal.Builder) {
	if len(j.childProcs) == 0 {
		return
	}

	for _, p := range j.childProcs {
		if err := p.Signal(os.Interrupt); err != nil {
			b.Logger.Error("SIGINT to child proc", err.Error())
		}
	}
}
