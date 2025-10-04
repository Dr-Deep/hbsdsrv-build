//go:build freebsd

package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	// Shell
	JobDefaultShell     = "/usr/local/bin/bash"
	JobDefaultShellArgs = "-euxc"

	// Niceness
	JobNiceness = 15

	// Proccess Shutdown Deadline
	JobDeadline = time.Second * 30
)

func (b *Builder) ShutdownProc(cmd *exec.Cmd, kill context.CancelFunc) error {
	var pid = cmd.Process.Pid

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		return err
	}

	b.Logger.Info("waiting for proccess shutdown", fmt.Sprintf("PID:%v", pid))

	// ticker
	var (
		deadline = time.After(JobDeadline)
		ticker   = time.NewTicker(500 * time.Millisecond)
	)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			kill()
			b.Logger.Error("killed because of deadline", fmt.Sprintf("PID:%v", pid))
			return nil

		case <-ticker.C:
			if cmd.Process != nil && cmd.ProcessState.Exited() {
				// process exited cleanly
				return nil
			}
		}
	}
}

// non-blocking; output goes to our stdin,stdout,stderr
func (b *Builder) RunOSCommand(argv []string) (*exec.Cmd, context.CancelFunc, error) {
	// Command Arg Vector
	var rawcmdargv = strings.Builder{}
	for _, s := range argv {
		rawcmdargv.WriteString(s)
		rawcmdargv.WriteRune(' ')
	}

	var cmdargv = []string{JobDefaultShellArgs, rawcmdargv.String()}

	// Command
	ctx, cancel := context.WithCancel(context.Background())
	var cmd = exec.CommandContext(
		ctx,
		JobDefaultShell,
		cmdargv...,
	)

	// reroute output
	if b.Logger.File == os.Stdout {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = b.Logger.File
	cmd.Stderr = b.Logger.File

	// Process Attributes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,           // Process Group
		Pdeathsig: syscall.SIGINT, // if we die
	}

	// Start Command
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, nil, err
	}

	// Niceness
	// https://man.freebsd.org/cgi/man.cgi?query=setpriority&manpath=OpenBSD+3.2
	if err := syscall.Setpriority(
		syscall.PRIO_PGRP, // job and all his childs
		cmd.Process.Pid,   // PID, group parent
		JobNiceness,
	); err != nil {
		cancel()
		return nil, nil, err
	}

	return cmd, cancel, nil
}
