package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/httpthings/run"
)

type Server struct {
	listener string
	base     string
	handler  http.Handler
	todos    TodoSet
}

func NewServer(listener, static string) (Server, error) {
	ts := NewTodoSet()
	// Populate some initial tasks
	ts.Add("Task1", "The Fist Task")
	ts.Add("Task2", "The Second Task")

	s := Server{
		listener: listener,
		todos:    *ts,
		base:     "api",
	}

	r := mux.NewRouter()

	routes := s.routes()
	err := routes.Populate(r, s.base)
	if err != nil {
		return s, err
	}

	if static != "" {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(static)))
	}

	s.handler = r
	return s, nil
}

func (s Server) run() {
	run.Run(run.DetectRunMode(), s.listener, s.handler, true)
}
