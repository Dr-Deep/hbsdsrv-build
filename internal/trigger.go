package internal

type Trigger interface {
	Run(b *Builder, c chan TriggerSignal)
}

type TriggerSignal struct {
	JobName string
	Reason  string
}

/*
triggers:

  - git-repo-url: https://git.hardenedbsd.org/hardenedbsd/HardenedBSD
    git-branch: hardened/14-stable/master
    job: pkgbase

  - git-repo-url: https://git.hardenedbsd.org/hardenedbsd/ports.git
    git-branch: hardenedbsd/main
    job: ports
*/
