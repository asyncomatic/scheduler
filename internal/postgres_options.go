package internal

import (
	"fmt"
	"github.com/caitlinelfring/go-env-default"
)

type PostgresOptions struct {
	host     string
	port     int
	username string
	password string
	database string
	sslmode  string
}

func NewPostgresOptions() *PostgresOptions {
	return &PostgresOptions{
		host:     env.GetDefault("POSTGRES_HOST", "localhost"),
		port:     env.GetIntDefault("POSTGRES_PORT", 5432),
		username: env.GetDefault("POSTGRES_USER", "postgres"),
		password: env.GetDefault("POSTGRES_PASSWORD", "password"),
		database: env.GetDefault("POSTGRES_DB", "devcloud"),
		sslmode:  env.GetDefault("POSTGRES_SSLMODE", "disable")}
}

func (o *PostgresOptions) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		o.host, o.port, o.username, o.password, o.database, o.sslmode)
}
