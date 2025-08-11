package jobs

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
)

type JobPorts struct {
	args config.JobArgs
}

func NewJobPorts(args config.JobArgs) internal.Job {
	return &JobPorts{args: args}
}

func (j *JobPorts) Run(b *internal.Builder) error {
	_, cmd := internal.RunCommand(j.args)
	cmd.Stdout = b.Logger.File
	cmd.Stderr = b.Logger.File

	return cmd.Run()
}

func (j *JobPorts) Abort() {}
