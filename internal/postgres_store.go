package internal

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"net/url"
	"strconv"
)

type PostgresStore struct {
	DbStore *sql.DB
}

func NewPostgresStore(opts *PostgresOptions) *PostgresStore {
	db, err := sql.Open("postgres", opts.ConnString())
	if err != nil {
		panic(err)
	}

	return &PostgresStore{db}
}

func (s PostgresStore) Add(job Job) (int, error) {
	id := 0
	sqlStatement := `INSERT INTO jobs (delay, queue, team_id, user_id, description, payload)
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := s.DbStore.QueryRow(sqlStatement,
		job.Delay, job.Queue, job.TeamId, job.UserId, job.Description, job.Payload).Scan(&id)

	return id, err
}

func (s PostgresStore) List(values url.Values) ([]Job, error) {
	jobs := make([]Job, 0)
	sqlStatement := `SELECT id, delay, queue, team_id, user_id, description, payload
					FROM jobs LIMIT 100`

	rows, err := s.DbStore.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		job := &Job{}
		err = rows.Scan(&job.Id, &job.Delay, &job.Queue, &job.TeamId,
			&job.UserId, &job.Description, &job.Payload)
		if err != nil {
			panic(err)
		}

		jobs = append(jobs, *job)
	}

	err = rows.Err()
	return jobs, err
}

func (s PostgresStore) Get(id int) (Job, error) {
	var job Job
	sqlStatement := `SELECT id, delay, queue, team_id, user_id, description, payload 
					FROM jobs WHERE id = $1`

	row := s.DbStore.QueryRow(sqlStatement, id)
	err := row.Scan(&job.Id, &job.Delay, &job.Queue, &job.TeamId,
		&job.UserId, &job.Description, &job.Payload)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return job, errors.New("Error getting job: " + strconv.Itoa(id) + " not found")
	case err == nil:
		return job, nil
	default:
		panic(err)
	}
}

func (s PostgresStore) Delete(id int) error {
	sqlStatement := `DELETE FROM jobs WHERE id = $1;`
	_, err := s.DbStore.Exec(sqlStatement, id)

	return err
}
