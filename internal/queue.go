package internal

type Queue interface {
	Write(job Job) error
}
