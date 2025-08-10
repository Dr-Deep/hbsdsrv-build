package jobs

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
)

type JobPkgbase struct{}

func NewJobPkgbase(config.JobArgs) internal.Job {
	return &JobPkgbase{}
}

func (j *JobPkgbase) Run(b *internal.Builder) error {
	return nil
}

func (j *JobPkgbase) Abort() {}
