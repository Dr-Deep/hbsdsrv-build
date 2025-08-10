package internal

import (
	"context"
	"os/exec"
	"strings"
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

	var sus = strings.Builder{}
	for _, s := range args {
		sus.WriteString(s)
		sus.WriteRune(' ')
	}

	var argv = []string{DefaultShellArgs, sus.String()}

	ctx, cancel := context.WithCancel(context.Background())
	var cmd = exec.CommandContext(
		ctx,
		DefaultShell,
		argv...,
	)

	return cancel, cmd
}
