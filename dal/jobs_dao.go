package dal

import (
	"database/sql"
	"reflect"
	"scheduler/dal/postgres"
	"scheduler/models"
)

var JobsDaoRegistry = map[string]JobsDao{
	"_default": postgres.NewPostgresJobsDao(),
	"postgres": postgres.NewPostgresJobsDao(),
}

type JobsDao interface {
	Add(*models.Job) (int, error)
	List() ([]models.Job, error)
	Get(int) (models.Job, error)
	Delete(int) error

	DBConn(*sql.DB)
}

func NewJobsDao(db *sql.DB, daoType string) JobsDao {
	dao := reflect.ValueOf(JobsDaoRegistry[daoType]).Interface().(JobsDao)
	dao.DBConn(db)

	return dao
}
