//go:build freebsd

package jobs

import (
	"context"

	"github.com/Dr-Deep/hbsdsrv-build/internal"
	"github.com/Dr-Deep/hbsdsrv-build/internal/config"
)

type JobCmd struct {
	args        config.JobArgs
	stopJobProc context.CancelFunc
}

func NewJobCmd(args config.JobArgs) internal.Job {
	return &JobCmd{args: args}
}

// run until command is done
func (j *JobCmd) Run(b *internal.Builder) error {
	cmd, cancel, err := b.RunOSCommand(j.args)
	if err != nil {
		return err
	}

	j.stopJobProc = cancel

	return cmd.Wait()
}

func (j *JobCmd) Abort(b *internal.Builder) error {
	if j.stopJobProc == nil {
		return nil
	}

	j.stopJobProc()

	return nil
}
