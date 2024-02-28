package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewPostgresDBConn() *sql.DB {
	db, err := sql.Open("postgres", NewPostgresOptions().ConnString())
	if err != nil {
		panic(err)
	}

	return db
}
