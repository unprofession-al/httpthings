package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/httpthings/respond"
)

func (s Server) OpenAPIHandler(res http.ResponseWriter, req *http.Request) {
	out, err := s.routes().AsHTML("test", s.base)
	if err != nil {
		respond.Auto(res, req, http.StatusInternalServerError, fmt.Sprintf("internal server error"))
		return
	}
	respond.Raw(res, http.StatusOK, out)
}

func (s Server) ListTodosHandler(res http.ResponseWriter, req *http.Request) {
	respond.Auto(res, req, http.StatusOK, s.todos.AsSlice())
}

func (s Server) ShowTodoHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name, ok := vars["todo"]
	if !ok {
		respond.Auto(res, req, http.StatusNotAcceptable, fmt.Sprintf("todo not provided"))
		return
	}
	if todo, found := s.todos[name]; found {
		respond.Auto(res, req, http.StatusOK, todo)
		return
	}
	respond.Auto(res, req, http.StatusNotFound, "not found")
}

func (s Server) AddTodoHandler(res http.ResponseWriter, req *http.Request) {
	todo := &Todo{}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respond.Auto(res, req, http.StatusInternalServerError, "could not read request body")
		return
	}
	err = json.Unmarshal(b, &todo)
	if err != nil {
		respond.Auto(res, req, http.StatusNotAcceptable, "could not unmarshal data")
		return
	}
	if _, found := s.todos[todo.Name]; found {
		respond.Auto(res, req, http.StatusConflict, todo)
		return
	}
	s.todos.Add(todo.Name, todo.Description)
	respond.Auto(res, req, http.StatusOK, todo)
}

func (s Server) FinishTodoHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name, ok := vars["todo"]
	if !ok {
		respond.Auto(res, req, http.StatusNotAcceptable, fmt.Sprintf("todo not provided"))
		return
	}
	if _, found := s.todos[name]; found {
		s.todos[name].Finish()
		respond.Auto(res, req, http.StatusOK, s.todos[name])
		return
	}
	respond.Auto(res, req, http.StatusNotFound, "not found")
}
