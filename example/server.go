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

func NewServer(listener, static string) Server {
	ts := NewTodoSet()
	ts.Add("Task1", "The Fist Task")
	ts.Add("Task2", "The Second Task")

	s := Server{
		listener: listener,
		todos:    *ts,
		base:     "api",
	}

	r := mux.NewRouter().StrictSlash(true)

	routes := s.routes()
	routes.Populate(r, s.base)

	if static != "" {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(static)))
	}

	s.handler = r
	return s
}

func (s Server) run() {
	run.Run(run.DetectRunMode(), s.listener, s.handler, true)
}
