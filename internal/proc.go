//go:build freebsd

package internal

import (
	"context"
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
)

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

	// chance to cleanup
	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGINT)
	}

	// reroute output
	if b.Logger.File == os.Stdout {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = b.Logger.File
	cmd.Stderr = b.Logger.File

	// kill delay
	cmd.WaitDelay = 30 * time.Second

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
