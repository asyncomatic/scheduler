//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

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
