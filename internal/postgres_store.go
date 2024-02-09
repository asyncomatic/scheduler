package internal

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"net/url"
	"os"
	"strconv"
)

type PostgresStore struct {
	DbStore *sql.DB
}

func NewPostgresStore() *PostgresStore {
	host := os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}
	username := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	//host := "localhost"
	//port := 5432
	//username := "postgres"
	//pass := "password"
	//dbname := "devcloud"

	options := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, pass, dbname)

	db, err := sql.Open("postgres", options)
	if err != nil {
		panic(err)
	}

	return &PostgresStore{db}
}

func (s PostgresStore) Add(test Test) (int, error) {
	id := 0
	sqlStatement := `
		INSERT INTO tests (delay, queue, team_id, user_id, description, payload)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := s.DbStore.QueryRow(sqlStatement,
		test.Delay, test.Queue, test.TeamId, test.UserId, test.Description, test.Payload).Scan(&id)

	return id, err
}

func (s PostgresStore) List(values url.Values) ([]Test, error) {
	tests := make([]Test, 0)
	sqlStatement := `
		SELECT id, delay, queue, team_id, user_id, description, payload 
		FROM tests LIMIT 100`

	rows, err := s.DbStore.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		test := &Test{}
		err = rows.Scan(&test.Id, &test.Delay, &test.Queue, &test.TeamId,
			&test.UserId, &test.Description, &test.Payload)
		if err != nil {
			panic(err)
		}

		tests = append(tests, *test)
	}

	err = rows.Err()
	return tests, err
}

func (s PostgresStore) Get(id int) (Test, error) {
	var test Test
	sqlStatement := `
		SELECT id, delay, queue, team_id, user_id, description, payload 
		FROM tests WHERE id = $1`

	row := s.DbStore.QueryRow(sqlStatement, id)
	err := row.Scan(&test.Id, &test.Delay, &test.Queue, &test.TeamId,
		&test.UserId, &test.Description, &test.Payload)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return test, errors.New("Error getting test: " + strconv.Itoa(id) + " not found")
	case err == nil:
		return test, nil
	default:
		panic(err)
	}
}

func (s PostgresStore) Delete(id int) error {
	sqlStatement := `DELETE FROM tests WHERE id = $1;`
	_, err := s.DbStore.Exec(sqlStatement, id)

	return err
}
