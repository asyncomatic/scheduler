//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

package postgres

import (
	"fmt"
	"github.com/caitlinelfring/go-env-default"
)

type Options struct {
	host     string
	port     int
	username string
	password string
	database string
	sslmode  string
}

func NewPostgresOptions() *Options {
	return &Options{
		host:     env.GetDefault("POSTGRES_HOST", "localhost"),
		port:     env.GetIntDefault("POSTGRES_PORT", 5432),
		username: env.GetDefault("POSTGRES_USER", "postgres"),
		password: env.GetDefault("POSTGRES_PASSWORD", "password"),
		database: env.GetDefault("POSTGRES_DB", "devcloud"),
		sslmode:  env.GetDefault("POSTGRES_SSLMODE", "disable")}
}

func (o *Options) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		o.host, o.port, o.username, o.password, o.database, o.sslmode)
}
