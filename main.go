//go:build freebsd

package main

import (
	"flag"
	"os"

	"github.com/Dr-Deep/hbsdsrv-build/internal"
	"github.com/Dr-Deep/hbsdsrv-build/internal/config"
	"github.com/Dr-Deep/hbsdsrv-build/internal/jobs"
	"github.com/Dr-Deep/hbsdsrv-build/internal/triggers"

	"github.com/Dr-Deep/logging-go"
)

const rwForOwnerOnlyPerm = 0o600

var (
	allTriggers = map[string]func(config.TriggerConfig) internal.Trigger{
		"git": triggers.NewTriggerGit,
	}

	allJobTargets = map[string]func(config.JobArgs) internal.Job{
		"hbsdsrv": jobs.NewJobCmd,
		"ports":   jobs.NewJobCmd,
		"pkgbase": jobs.NewJobPkgbase,
	}

	cfg    *config.Configuration
	logger *logging.Logger

	// Flags
	configFilePath = flag.String(
		"config-file",
		"./config.yml",
		"configuration file",
	)
	loggingFilePath = flag.String(
		"logging-file",
		"",
		"logging file",
	)
)

func init() {
	flag.Parse()
}

func setupConfig() {
	// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
	cfgFile, err := os.OpenFile(
		*configFilePath,
		os.O_RDONLY,
		rwForOwnerOnlyPerm,
	)
	if err != nil {
		panic(err)
	}

	info, err := cfgFile.Stat()
	if err != nil {
		panic(err)
	}

	// check if the config file is ro
	if info.Mode().Perm()&0o022 != 0 {
		panic("config file must not be writable")
	}

	cfg, err = config.LoadConfig(cfgFile)
	if err != nil {
		panic(err)
	}
}

func setupLogger() {
	var logFile *os.File

	// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
	if *loggingFilePath != "" {
		_logFile, err := os.OpenFile(
			*loggingFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			rwForOwnerOnlyPerm,
		)
		if err != nil {
			panic(err)
		}
		logFile = _logFile
	}

	logger = logging.NewLogger(logFile)
}

func main() {
	setupConfig()
	setupLogger()

	defer func() {
		if err := logger.Close(); err != nil {
			panic(err)
		}
	}()

	var builder = internal.NewBuilder(logger, cfg)

	// Register Triggers
	for _, t := range cfg.Triggers {
		for s, f := range allTriggers {
			if s == t.Type {
				builder.RegisterTrigger(f, t)
			}
		}
	}

	// Register Jobs
	for job, targets := range cfg.Jobs {
		for target, targetArgs := range targets {
			for _, arg := range targetArgs {
				f, oke := allJobTargets[job]
				if oke {
					builder.RegisterJob(
						job,
						target,
						f(arg),
					)
				} else {
					logger.Error("Job not registered", job, target)
				}
			}
		}
	}

	builder.Launch()
}
