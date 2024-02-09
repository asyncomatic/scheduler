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

type TestsHandler struct {
	store Store
	queue Queue
}

var (
	routeRegex       = regexp.MustCompile(`^(/tests)/*$`)
	routeWithIdRegex = regexp.MustCompile(`^/tests/([a-z0-9]+(?:-[a-z0-9]+)*)/*$`)

	tests = make(map[string]string)
)

func NewTestsHandler(store Store, queue Queue) *TestsHandler {
	return &TestsHandler{
		store: store,
		queue: queue,
	}
}

func (h *TestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && routeRegex.MatchString(r.URL.Path):
		h.create(w, r)
		return

	case r.Method == http.MethodGet && routeRegex.MatchString(r.URL.Path):
		h.list(w, r)
		return

	case r.Method == http.MethodGet && routeWithIdRegex.MatchString(r.URL.Path):
		h.get(w, r)
		return

	case r.Method == http.MethodDelete && routeWithIdRegex.MatchString(r.URL.Path):
		h.delete(w, r)
		return

	default:
		onError(w, http.StatusNotImplemented)
		return
	}
}

func (h *TestsHandler) create(w http.ResponseWriter, r *http.Request) {
	var test Test

	err := json.NewDecoder(r.Body).Decode(&test)
	if err != nil {
		onError(w, http.StatusBadRequest, "Invalid format for test")
		return
	}

	id, err := h.store.Add(test)
	if err != nil {
		onError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func(id int, delay int) {
		time.Sleep(time.Duration(delay) * time.Second)
		t, err := h.store.Get(id)
		if err != nil {
			fmt.Println("Test " + strconv.Itoa(id) + " not found, skipping execution")
			return
		}

		fmt.Println("Test " + strconv.Itoa(id) + " enqueued for execution")
		_ = h.queue.Write(t)
		b, err := json.Marshal(t)
		fmt.Println(string(b))
		_ = h.store.Delete(id)
		return

	}(id, test.Delay)

	w.Header().Set("Location", r.Host+r.RequestURI+"/"+strconv.Itoa(id))
	w.WriteHeader(http.StatusCreated)

	return
}

func (h *TestsHandler) list(w http.ResponseWriter, r *http.Request) {
	tests, err := h.store.List(r.URL.Query())
	if err != nil {
		onError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tests)
	if err != nil {
		onError(w, http.StatusInternalServerError)
		return
	}

	return
}

func (h *TestsHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		onError(w, http.StatusBadRequest, "Missing or invalid test id '"+
			routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1]+
			"'")
		return
	}

	test, err := h.store.Get(id)
	if err != nil {
		onError(w, http.StatusNotFound, "Test with id '"+
			routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1]+
			"' not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(test)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func (h *TestsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1])
	if err != nil {
		onError(w, http.StatusBadRequest, "Missing or invalid test id '"+
			routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1]+
			"'")
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		onError(w, http.StatusNotFound, "Test with id '"+
			routeWithIdRegex.FindStringSubmatch(r.URL.Path)[1]+
			"' not found")
		return
	}

	w.Header().Set("Location", r.Host+r.RequestURI)
	w.WriteHeader(http.StatusAccepted)
	return
}

func onError(w http.ResponseWriter, status int, msg ...string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(strings.Join(msg, " ")))
	if err != nil {

	}
}
