package internal

type Queue interface {
	Write(test Test) error
}
