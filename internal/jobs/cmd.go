//go:build freebsd

package jobs

import (
	"context"
	"os/exec"

	"github.com/Dr-Deep/hbsdsrv-build/internal"
	"github.com/Dr-Deep/hbsdsrv-build/internal/config"
)

type JobCmd struct {
	args config.JobArgs

	jobProc     *exec.Cmd
	killJobProc context.CancelFunc
}

func NewJobCmd(args config.JobArgs) internal.Job {
	return &JobCmd{args: args}
}

// run until command is done
func (j *JobCmd) Run(b *internal.Builder) error {
	cmd, kill, err := b.RunOSCommand(j.args)
	if err != nil {
		return err
	}

	j.jobProc = cmd
	j.killJobProc = kill

	return cmd.Wait()
}

func (j *JobCmd) Abort(b *internal.Builder) error {
	if j.jobProc == nil {
		return nil
	}

	return b.ShutdownProc(j.jobProc, j.killJobProc)
}
