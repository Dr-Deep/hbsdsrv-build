package internal

import (
	"context"
	"os/exec"
)

const DefaultShell = "/bin/sh"
const DefaultShellArgs = "-euxc"

type Job interface {
	Run(*Builder) error
	Abort()
}

func RunCommand(args []string) (context.CancelFunc, *exec.Cmd) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	var argv = []string{DefaultShellArgs}
	argv = append(argv, args...)

	ctx, cancel := context.WithCancel(context.Background())
	var cmd = exec.CommandContext(
		ctx,
		DefaultShell,
		argv...,
	)

	return cancel, cmd
}
