package models

type Job struct {
	Id         int    `json:"id"`
	Delay      int    `json:"delay"`
	Queue      string `json:"queue"`
	Class      string `json:"class"`
	Method     string `json:"method"`
	RetryCount int    `json:"retry_count"`
	State      string `json:"state"`
}
