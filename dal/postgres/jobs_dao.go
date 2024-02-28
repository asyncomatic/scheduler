package postgres

import (
	"database/sql"
	"errors"
	"github.com/caitlinelfring/go-env-default"
	_ "github.com/lib/pq"
	"scheduler/models"
	"strconv"
)

type JobsDao struct {
	db *sql.DB
}

func NewPostgresJobsDao() *JobsDao {
	return &JobsDao{}
}

func (s *JobsDao) DBConn(db *sql.DB) {
	s.db = db
}

func (s *JobsDao) Add(job *models.Job) (int, error) {
	id := 0
	sqlStatement := `INSERT INTO jobs (delay, queue, team_id, user_id, description, payload)
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := s.db.QueryRow(sqlStatement,
		job.Delay, job.Queue, job.TeamId, job.UserId, job.Description, job.Payload).Scan(&id)

	return id, err
}

func (s *JobsDao) List() ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	sqlStatement := `SELECT id, delay, queue, team_id, user_id, description, payload
					FROM jobs LIMIT 100`

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		job := &models.Job{}
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

func (s *JobsDao) Get(id int) (models.Job, error) {
	var job models.Job
	sqlStatement := `SELECT id, delay, queue, team_id, user_id, description, payload 
					FROM jobs WHERE id = $1`

	row := s.db.QueryRow(sqlStatement, id)
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

func (s *JobsDao) Delete(id int) error {
	sqlStatement := `DELETE FROM jobs WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, id)

	return err
}

func (s *JobsDao) hydrate(token string, job *models.Job) error {
	job.Queue = env.GetDefault("KAFKA_DEFAULT_TOPIC", "devcloud")

	return nil
}