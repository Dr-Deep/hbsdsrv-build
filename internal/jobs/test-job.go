package jobs

import (
	"fmt"
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
)

type JobTest struct {
	args []string
}

func NewJobTest(av config.JobArgs) internal.Job {
	return &JobTest{
		args: av,
	}
}

func (j *JobTest) Run(b *internal.Builder) error {
	b.Logger.Info(
		fmt.Sprintf("running: %v", j.args),
	)

	ctx, cmd := internal.RunCommand(j.args)
	cmd.Stdout = b.Logger.File
	cmd.Stderr = b.Logger.File

	b.RunningJobs = append(b.RunningJobs, ctx)

	return cmd.Run()
}

func (j JobTest) Abort() {}
