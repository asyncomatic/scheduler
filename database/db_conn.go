//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

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
