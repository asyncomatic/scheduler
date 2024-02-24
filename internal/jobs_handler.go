package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JobsHandler struct {
	regex   *regexp.Regexp
	regexId *regexp.Regexp
	store   Store
	queue   Queue
}

func NewJobsHandler(route string, store Store, queue Queue) *JobsHandler {
	return &JobsHandler{
		regex:   regexp.MustCompile(`^(` + route + `)/*$`),
		regexId: regexp.MustCompile(`^` + route + `/([a-z0-9]+(?:-[a-z0-9]+)*)/*$`),
		store:   store,
		queue:   queue,
	}
}

func (h *JobsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && h.regex.MatchString(r.URL.Path):
		h.create(w, r)
		return

	case r.Method == http.MethodGet && h.regex.MatchString(r.URL.Path):
		h.list(w, r)
		return

	case r.Method == http.MethodGet && h.regexId.MatchString(r.URL.Path):
		h.get(w, r)
		return

	case r.Method == http.MethodDelete && h.regexId.MatchString(r.URL.Path):
		h.delete(w, r)
		return

	default:
		msg := fmt.Sprintf("Route '%s' not supprted", r.URL.Path)
		reply(w, http.StatusNotImplemented, msg)
		return
	}
}

func (h *JobsHandler) create(w http.ResponseWriter, r *http.Request) {
	var job Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		msg := fmt.Sprintf("Invalid json format for job: %s", err.Error())
		reply(w, http.StatusBadRequest, msg)
		return
	}

	id, err := h.store.Add(job)
	if err != nil {
		msg := fmt.Sprintf("Error storing job: %s", err.Error())
		reply(w, http.StatusInternalServerError, msg)
		return
	}

	go func(id int, delay int) {
		time.Sleep(time.Duration(delay) * time.Second)
		t, err := h.store.Get(id)
		if err != nil {
			fmt.Println("Job " + strconv.Itoa(id) + " not found, skipping execution")
			return
		}

		fmt.Println("Job " + strconv.Itoa(id) + " enqueued for execution")
		_ = h.queue.Write(t)
		b, err := json.Marshal(t)
		fmt.Println(string(b))
		_ = h.store.Delete(id)
		return

	}(id, job.Delay)

	w.Header().Set("Location", r.Host+r.RequestURI+"/"+strconv.Itoa(id))
	w.WriteHeader(http.StatusCreated)

	return
}

func (h *JobsHandler) list(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.store.List(r.URL.Query())
	if err != nil {
		msg := fmt.Sprintf("Error retrieving list of jobs: %s", err.Error())
		reply(w, http.StatusInternalServerError, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jobs)
	if err != nil {
		msg := fmt.Sprintf("Error encoding list of jobs: %s", err.Error())
		reply(w, http.StatusInternalServerError, msg)
		return
	}

	return
}

func (h *JobsHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(h.regexId.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		msg := fmt.Sprintf("Missing or invalid format for job id '%s'",
			h.regexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusBadRequest, msg)
		return
	}

	job, err := h.store.Get(id)
	if err != nil {
		msg := fmt.Sprintf("Job with id '%s' not found",
			h.regexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusNotFound, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		msg := fmt.Sprintf("Error encoding job: %s", err.Error())
		reply(w, http.StatusInternalServerError, msg)
		return
	}

	return
}

func (h *JobsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(h.regexId.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		msg := fmt.Sprintf("Missing or invalid format for job id '%s'",
			h.regexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusBadRequest, msg)
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		msg := fmt.Sprintf("Job with id '%s' not found",
			h.regexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusNotFound, msg)
		return
	}

	w.Header().Set("Location", r.Host+r.RequestURI)
	reply(w, http.StatusAccepted)
	return
}

func reply(w http.ResponseWriter, status int, msg ...string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(strings.Join(msg, " ")))
	if err != nil {

	}
}
