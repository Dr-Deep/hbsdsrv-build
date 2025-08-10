package internal

import (
	"context"
	"os/exec"
)

const DefaultShell = "/bin/sh -euxc"

type Job interface {
	Run(*Builder) error
	Abort()
}

func RunCommand(args []string) (context.CancelFunc, *exec.Cmd) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	ctx, cancel := context.WithCancel(context.Background())
	var cmd = exec.CommandContext(
		ctx,
		DefaultShell,
		args...,
	)

	return cancel, cmd
}
