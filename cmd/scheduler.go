package main

import (
	"database/sql"
	"fmt"
	"github.com/caitlinelfring/go-env-default"
	"net/http"
	"scheduler/internal"
)

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	route := env.GetDefault("SCHEDULER_ROUTE", "/jobs")
	port := fmt.Sprintf(":%s", env.GetDefault("SCHEDULER_PORT", "8080"))

	storeopts := internal.NewPostgresOptions()
	store := internal.NewPostgresStore(storeopts)
	defer func(DbStore *sql.DB) {
		_ = DbStore.Close()
	}(store.DbStore)

	queueopts := internal.NewKafkaOptions()
	queue := internal.NewKafkaQueue(queueopts)

	jobsHandler := internal.NewJobsHandler(route, store, queue)

	mux := http.NewServeMux()

	mux.Handle(route, jobsHandler)
	mux.Handle(fmt.Sprintf("%s/", route), jobsHandler)

	mux.Handle("/", http.HandlerFunc(routeNotFound))

	err := http.ListenAndServe(port, mux)
	if err != nil {
		return
	}
}
