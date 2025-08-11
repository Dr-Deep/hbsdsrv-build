package main

import (
	"flag"
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"hbsdsrv-build/internal/jobs"
	"hbsdsrv-build/internal/logging"
	"hbsdsrv-build/internal/triggers"
	"os"
)

var (
	allTriggers = map[string]func(config.TriggerConfig) internal.Trigger{
		"git":  triggers.NewTriggerGit,
		"test": triggers.NewTriggerTest,
	}

	allJobTargets = map[string]func(config.JobArgs) internal.Job{
		"test":    jobs.NewJobTest,
		"ports":   jobs.NewJobPorts,
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

	// Todo: gucken das nur root schreibrechte auf cfg hat sonst führen wir hier ganix aus
}

func setupConfig() {
	// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
	cfgFile, err := os.OpenFile(
		*configFilePath,
		os.O_RDONLY,
		0600,
	)
	if err != nil {
		panic(err)
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
			0600,
		)
		if err != nil {
			panic(err)
		}
		logFile = _logFile
	} else {
		logFile = os.Stdout
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
