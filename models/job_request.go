package models

type JobRequest struct {
	Delay       int    `json:"delay"`
	Description string `json:"description"`
	Payload     string `json:"payload"`
}
