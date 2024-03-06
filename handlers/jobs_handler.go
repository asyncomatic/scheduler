package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"scheduler/dal"
	"scheduler/models"
	"scheduler/queue"
	"strconv"
	"strings"
	"time"
)

var jobsRegex = regexp.MustCompile(`^(/jobs)/*$`)
var jobsRegexId = regexp.MustCompile(`^/jobs/([a-z0-9]+(?:-[a-z0-9]+)*)/*$`)

type JobsHandler struct {
	jobsDao dal.JobsDao
	writer  queue.QueueWriter
}

func NewJobsHandler(jobsDao dal.JobsDao, writer queue.QueueWriter) *JobsHandler {
	return &JobsHandler{
		jobsDao: jobsDao,
		writer:  writer,
	}
}

func (h *JobsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && jobsRegex.MatchString(r.URL.Path):
		h.list(w, r)
		return

	case r.Method == http.MethodPost && jobsRegex.MatchString(r.URL.Path):
		h.create(w, r)
		return

	case r.Method == http.MethodGet && jobsRegexId.MatchString(r.URL.Path):
		h.read(w, r)
		return

	case r.Method == http.MethodDelete && jobsRegexId.MatchString(r.URL.Path):
		h.delete(w, r)
		return

	default:
		msg := fmt.Sprintf("Route '%s' not supprted", r.URL.Path)
		reply(w, http.StatusNotImplemented, msg)
		return
	}
}

func (h *JobsHandler) list(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobsDao.List()
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

func (h *JobsHandler) create(w http.ResponseWriter, r *http.Request) {
	var jobReq models.JobRequest

	err := json.NewDecoder(r.Body).Decode(&jobReq)
	if err != nil {
		msg := fmt.Sprintf("Invalid json format for job: %s", err.Error())
		reply(w, http.StatusBadRequest, msg)
		return
	}

	job := models.Job{
		Delay:      jobReq.Delay,
		Queue:      jobReq.Queue,
		Class:      jobReq.Class,
		Method:     jobReq.Method,
		RetryCount: jobReq.RetryCount,
	}

	if jobReq.State == nil {
		job.State = "{}"
	} else {
		state, err := json.Marshal(jobReq.State)
		if err != nil {
			msg := fmt.Sprintf("Error marshalling job state: %s", err.Error())
			reply(w, http.StatusInternalServerError, msg)
			return
		}
		job.State = string(state)
	}

	id, err := h.jobsDao.Add(&job)
	if err != nil {
		msg := fmt.Sprintf("Error storing job: %s", err.Error())
		reply(w, http.StatusInternalServerError, msg)
		return
	}

	go func(id int, delay int) {
		time.Sleep(time.Duration(delay) * time.Second)
		t, err := h.jobsDao.Get(id)
		if err != nil {
			fmt.Println("Job " + strconv.Itoa(id) + " not found, skipping execution")
			return
		}

		fmt.Println("Job " + strconv.Itoa(id) + " enqueued for execution")
		_ = h.writer.Write(t)
		b, err := json.Marshal(t)
		fmt.Println(string(b))
		_ = h.jobsDao.Delete(id)
		return

	}(id, job.Delay)

	w.Header().Set("Location", r.Host+r.RequestURI+"/"+strconv.Itoa(id))
	w.WriteHeader(http.StatusCreated)

	return
}

func (h *JobsHandler) read(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(jobsRegexId.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		msg := fmt.Sprintf("Missing or invalid format for job id '%s'",
			jobsRegexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusBadRequest, msg)
		return
	}

	job, err := h.jobsDao.Get(id)
	if err != nil {
		msg := fmt.Sprintf("Job with id '%s' not found",
			jobsRegexId.FindStringSubmatch(r.URL.Path)[1])

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
	id, err := strconv.Atoi(jobsRegexId.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		msg := fmt.Sprintf("Missing or invalid format for job id '%s'",
			jobsRegexId.FindStringSubmatch(r.URL.Path)[1])

		reply(w, http.StatusBadRequest, msg)
		return
	}

	err = h.jobsDao.Delete(id)
	if err != nil {
		msg := fmt.Sprintf("Job with id '%s' not found",
			jobsRegexId.FindStringSubmatch(r.URL.Path)[1])

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
