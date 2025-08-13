package internal

import (
	"context"
	"os/exec"
	"strings"
)

const DefaultShell = "/bin/sh"
const DefaultShellArgs = "-euxc"

type Job interface {
	Run(b *Builder) error
	Abort(b *Builder)
}

/*
	// Befehl vorbereiten
	cmd := exec.Command("sleep", "100")

	// Neue Prozessgruppe erstellen
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // eigene Prozessgruppe
	}

	// Starten
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Etwas warten
	time.Sleep(2 * time.Second)

	// SIGINT an gesamte Prozessgruppe senden
	// Negative PID => Prozessgruppe statt einzelner Prozess
	syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)

	// Warten bis alles beendet ist
	cmd.Wait()
*/

func RunCommand(args []string) (context.CancelFunc, *exec.Cmd) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	var cmdargv = strings.Builder{}
	for _, s := range args {
		cmdargv.WriteString(s)
		cmdargv.WriteRune(' ')
	}

	var argv = []string{DefaultShellArgs, cmdargv.String()}

	ctx, cancel := context.WithCancel(context.Background())
	var cmd = exec.CommandContext(
		ctx,
		DefaultShell,
		argv...,
	)

	return cancel, cmd
}
