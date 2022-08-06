package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s Server) OpenAPIHandler(res http.ResponseWriter, req *http.Request) {
	out, err := s.routes().AsHTML("test", s.base)
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, fmt.Sprintf("internal server error"))
		return
	}
	s.raw(res, http.StatusOK, out)
}

func (s Server) ListTodosHandler(res http.ResponseWriter, req *http.Request) {
	s.respond(res, req, http.StatusOK, s.todos)
}

func (s Server) ShowTodoHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name, ok := vars["todo"]
	if !ok {
		s.respond(res, req, http.StatusNotFound, fmt.Sprintf("todo not provided"))
		return
	}
	if todo, found := s.todos[name]; found {
		s.respond(res, req, http.StatusOK, todo)
		return
	}
	s.respond(res, req, http.StatusNotFound, "not found")
}

func (s Server) AddTodoHandler(res http.ResponseWriter, req *http.Request) {
	s.respond(res, req, http.StatusNotImplemented, "not implemented")
}

func (s Server) FinishTodoHandler(res http.ResponseWriter, req *http.Request) {
	s.respond(res, req, http.StatusNotImplemented, "not implemented")
}
