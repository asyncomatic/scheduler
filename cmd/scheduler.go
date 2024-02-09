package main

import (
	"database/sql"
	"net/http"
	"scheduler/internal"
)

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	store := internal.NewPostgresStore()
	defer func(DbStore *sql.DB) {
		_ = DbStore.Close()
	}(store.DbStore)

	queue := internal.NewKafkaQueue()
	testsHandler := internal.NewTestsHandler(store, queue)

	mux := http.NewServeMux()

	mux.Handle("/tests", testsHandler)
	mux.Handle("/tests/", testsHandler)

	mux.Handle("/", http.HandlerFunc(routeNotFound))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
