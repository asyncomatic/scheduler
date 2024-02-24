package internal

import "net/url"

type Store interface {
	Add(job Job) (int, error)
	List(values url.Values) ([]Job, error)
	Get(id int) (Job, error)
	Delete(id int) error
}
