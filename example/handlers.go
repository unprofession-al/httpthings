package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/httpthings/respond"
)

func (s Server) OpenAPIHandler(res http.ResponseWriter, req *http.Request) {
	out, err := s.routes().AsHTML("test", s.base)
	if err != nil {
		respond.Respond(res, req, http.StatusInternalServerError, fmt.Sprintf("internal server error"))
		return
	}
	respond.Raw(res, http.StatusOK, out)
}

func (s Server) ListTodosHandler(res http.ResponseWriter, req *http.Request) {
	respond.Respond(res, req, http.StatusOK, s.todos)
}

func (s Server) ShowTodoHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name, ok := vars["todo"]
	if !ok {
		respond.Respond(res, req, http.StatusNotFound, fmt.Sprintf("todo not provided"))
		return
	}
	if todo, found := s.todos[name]; found {
		respond.Respond(res, req, http.StatusOK, todo)
		return
	}
	respond.Respond(res, req, http.StatusNotFound, "not found")
}

func (s Server) AddTodoHandler(res http.ResponseWriter, req *http.Request) {
	respond.Respond(res, req, http.StatusNotImplemented, "not implemented")
}

func (s Server) FinishTodoHandler(res http.ResponseWriter, req *http.Request) {
	respond.Respond(res, req, http.StatusNotImplemented, "not implemented")
}
