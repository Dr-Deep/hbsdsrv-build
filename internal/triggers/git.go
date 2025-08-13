package triggers

import (
	"hbsdsrv-build/internal"
	"hbsdsrv-build/internal/config"
	"time"

	git "github.com/go-git/go-git/v6"
	gitcfg "github.com/go-git/go-git/v6/config"
	gitplumbing "github.com/go-git/go-git/v6/plumbing"
	gitmem "github.com/go-git/go-git/v6/storage/memory"
)

// alle paar min (15/30min) abfragen
const gitIntervalDuration = time.Minute * 15

type TriggerGit struct {
	remote         *git.Remote
	lastCommitHash gitplumbing.Hash

	config.TriggerConfig
}

func NewTriggerGit(cfg config.TriggerConfig) internal.Trigger {
	var self = &TriggerGit{
		TriggerConfig: cfg,
	}

	// remote in memory
	self.remote = git.NewRemote(
		gitmem.NewStorage(),
		&gitcfg.RemoteConfig{
			Name: "origin",
			URLs: []string{cfg.GitRepoURL},
		})

	// refs
	refs, err := self.remote.List(&git.ListOptions{})
	if err != nil {
		panic(err)
	}

	// get latest hash
	var (
		branch  = cfg.GitBranch
		refName = gitplumbing.NewBranchReferenceName(branch)
	)

	for _, r := range refs {
		if r.Name() == refName {
			self.lastCommitHash = r.Hash()
		}
	}

	return self
}

func (t *TriggerGit) Run(b *internal.Builder, c chan internal.TriggerSignal) {
	var checkForNewCommits = func() {
		// get references
		refs, err := t.remote.List(&git.ListOptions{})
		if err != nil {
			b.Logger.Error("Git", err.Error())

			return
		}

		// latest hash?
		for _, r := range refs {
			if r.Name() == gitplumbing.NewBranchReferenceName(t.GitBranch) {
				// compare hashes from branch
				if r.Hash().String() != t.lastCommitHash.String() {
					// new commit
					t.lastCommitHash = r.Hash()

					b.Logger.Info("Git", t.Job, "got new commit", r.Hash().String(), r.String())

					// job losschicken
					c <- internal.TriggerSignal{
						JobName: t.Job,
						Reason:  "New Commit: " + r.String(),
					}

					return
				}

				b.Logger.Info("Git", t.Job, "were on head", t.lastCommitHash.String())
			}
		}
	}

	// trigger loop
	checkForNewCommits()

	// job losschicken
	c <- internal.TriggerSignal{
		JobName: t.Job,
		Reason:  "first run",
	}

	for {
		time.Sleep(gitIntervalDuration)
		checkForNewCommits()
	}
}
