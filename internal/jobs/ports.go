package jobs

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"os"
)

type JobPorts struct {
	args       config.JobArgs
	childProcs []*os.Process
}

func NewJobPorts(args config.JobArgs) internal.Job {
	return &JobPorts{args: args}
}

func (j *JobPorts) Run(b *internal.Builder) error {
	_, cmd := internal.RunCommand(j.args)
	cmd.Stdout = b.Logger.File
	cmd.Stderr = b.Logger.File

	j.childProcs = append(j.childProcs, cmd.Process)

	return cmd.Run()
}

func (j *JobPorts) Abort(b *internal.Builder) {
	if len(j.childProcs) == 0 {
		return
	}

	for _, p := range j.childProcs {
		if err := p.Signal(os.Interrupt); err != nil {
			b.Logger.Error("SIGINT to child proc", err.Error())
		}
	}
}
