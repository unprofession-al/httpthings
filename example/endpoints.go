package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/unprofession-al/httpthings/endpoint"
	"github.com/unprofession-al/httpthings/respond"
)

func (s Server) ListTodoEndpoint() *endpoint.Endpoint {
	ep := &endpoint.Endpoint{}
	ep.Name = "list-todos"
	ep.Responses = map[int]interface{}{http.StatusOK: []Todo{}}
	ep.ErrorResponse = HTTPError{}
	ep.Auth = s.auth
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		respond.Auto(w, r, http.StatusOK, s.todos.AsSlice())
	}
	return ep
}

func (s Server) ShowTodoEndpoint() *endpoint.Endpoint {
	ep := &endpoint.Endpoint{}
	ep.Name = "show-todo"
	ep.Responses = map[int]interface{}{http.StatusOK: Todo{}}
	ep.ErrorResponse = HTTPError{}
	errTodoNotProvided := ep.RegisterError(http.StatusNotAcceptable, "todo not provided")
	errTodoNotFound := ep.RegisterError(http.StatusNotFound, "todo not found")
	ep.Auth = s.auth
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		name, ok := ep.GetParamAsString("name", r)
		if !ok || len(name) < 0 {
			errTodoNotProvided(w, r)
			return
		}
		if todo, found := s.todos[name]; found {
			respond.Auto(w, r, http.StatusOK, todo)
			return
		}
		errTodoNotFound(w, r)
	}
	return ep
}

func (s Server) AddTodoEndpoint() *endpoint.Endpoint {
	ep := &endpoint.Endpoint{}
	ep.Name = "add-todo"
	ep.RequestBody = Todo{}
	ep.Responses = map[int]interface{}{http.StatusOK: Todo{}}
	ep.ErrorResponse = HTTPError{}
	errCouldNotReadRequest := ep.RegisterError(http.StatusInternalServerError, "could not read request")
	errCouldNotParseData := ep.RegisterError(http.StatusNotAcceptable, "could not parse data")
	errAlreadyExists := ep.RegisterError(http.StatusConflict, "todo already exists")
	ep.Auth = s.auth
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		todo := &TodoRequest{}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			errCouldNotReadRequest(w, r)
			return
		}
		err = json.Unmarshal(b, &todo)
		if err != nil {
			errCouldNotParseData(w, r)
			return
		}
		if _, found := s.todos[todo.Name]; found {
			errAlreadyExists(w, r)
			return
		}
		s.todos.Add(todo.AsTodo())
		respond.Auto(w, r, http.StatusOK, todo)
	}
	return ep
}

func (s Server) FinishTodoEndpoint() *endpoint.Endpoint {
	ep := &endpoint.Endpoint{}
	ep.Name = "finish-todo"
	ep.Responses = map[int]interface{}{http.StatusOK: Todo{}}
	ep.ErrorResponse = HTTPError{}
	errTodoNotProvided := ep.RegisterError(http.StatusNotAcceptable, "todo not provided")
	errTodoNotFound := ep.RegisterError(http.StatusNotFound, "todo not found")
	ep.Auth = s.auth
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		name, ok := ep.GetParamAsString("name", r)
		if !ok || len(name) < 1 {
			errTodoNotProvided(w, r)
			return
		}
		if _, found := s.todos[name]; found {
			s.todos[name].Finish()
			respond.Auto(w, r, http.StatusOK, s.todos[name])
			return
		}
		errTodoNotFound(w, r)
	}
	return ep
}
