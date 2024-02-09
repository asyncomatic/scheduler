package internal

import "net/url"

type Store interface {
	Add(test Test) (int, error)
	List(values url.Values) ([]Test, error)
	Get(id int) (Test, error)
	Delete(id int) error
}
