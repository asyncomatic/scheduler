package database

import (
	"database/sql"
	"reflect"
	"scheduler/database/postgres"
)

var DBConnRegistry = map[string]*sql.DB{
	"_default": postgres.NewPostgresDBConn(),
	"postgres": postgres.NewPostgresDBConn(),
}

func NewDBConn(dbType string) *sql.DB {
	return reflect.ValueOf(DBConnRegistry[dbType]).Interface().(*sql.DB)
}
