package internal

type Job interface {
	Run(b *Builder) error
	Abort(b *Builder) error
}
