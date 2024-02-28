package main

import (
	"database/sql"
	"fmt"
	"github.com/caitlinelfring/go-env-default"
	"net/http"
	"scheduler/auth"
	"scheduler/dal"
	"scheduler/database"
	"scheduler/handlers"
	"scheduler/queue"
)

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	queueType := env.GetDefault("SCHEDULER_QUEUE_TYPE", "_default")
	queueWriter := queue.NewQueueWriter(queueType)

	dbType := env.GetDefault("SCHEDULER_DB_TYPE", "_default")
	db := database.NewDBConn(dbType)
	defer func(conn *sql.DB) {
		_ = conn.Close()
	}(db)

	authType := env.GetDefault("SCHEDULER_AUTH_TYPE", "_default")
	authHandler := auth.NewAuthHandler(authType)

	jobsHandler := handlers.NewJobsHandler(dal.NewJobsDao(db, dbType), queueWriter)

	mux := http.NewServeMux()
	mux.Handle("/jobs", authHandler.Authn(jobsHandler.ServeHTTP))
	mux.Handle("/jobs/", authHandler.Authn(jobsHandler.ServeHTTP))

	mux.Handle("/", http.HandlerFunc(routeNotFound))

	port := fmt.Sprintf(":%s", env.GetDefault("SCHEDULER_PORT", "8080"))
	err := http.ListenAndServe(port, mux)
	if err != nil {
		return
	}
}
