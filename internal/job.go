package internal

type Job struct {
	Id          int    `json:"id"`
	Delay       int    `json:"delay"`
	Queue       string `json:"queue"`
	TeamId      int    `json:"team_id"`
	UserId      int    `json:"user_id"`
	Description string `json:"description"`
	Payload     string `json:"payload"`
}
