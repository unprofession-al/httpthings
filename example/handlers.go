package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/unprofession-al/httpthings/respond"
	"github.com/unprofession-al/httpthings/route"
)

func (s Server) ListTodosHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		respond.Auto(res, req, http.StatusOK, s.todos.AsSlice())
	}
}

func (s Server) ShowTodoHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params, errs := e.GetParams(req)
		if len(errs) > 0 {
			eOut := []string{}
			for _, e := range errs {
				eOut = append(eOut, e.Error())
			}
			respond.Auto(res, req, http.StatusNotAcceptable, eOut)
			return
		}
		name, ok := params["name"]
		if !ok || len(name) > 1 {
			respond.Auto(res, req, http.StatusNotAcceptable, "todo not provided")
			return
		}
		if todo, found := s.todos[name[0]]; found {
			respond.Auto(res, req, http.StatusOK, todo)
			return
		}
		respond.Auto(res, req, http.StatusNotFound, "not found")
	}
}

func (s Server) AddTodoHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		todo := &TodoRequest{}
		b, err := io.ReadAll(req.Body)
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
		s.todos.Add(todo.AsTodo())
		respond.Auto(res, req, http.StatusOK, todo)
	}
}

func (s Server) FinishTodoHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params, errs := e.GetParams(req)
		if len(errs) > 0 {
			eOut := []string{}
			for _, e := range errs {
				eOut = append(eOut, e.Error())
			}
			respond.Auto(res, req, http.StatusNotAcceptable, eOut)
			return
		}
		name, ok := params["name"]
		if !ok || len(name) > 1 {
			respond.Auto(res, req, http.StatusNotAcceptable, "todo not provided")
			return
		}
		if _, found := s.todos[name[0]]; found {
			s.todos[name[0]].Finish()
			respond.Auto(res, req, http.StatusOK, s.todos[name[0]])
			return
		}
		respond.Auto(res, req, http.StatusNotFound, "not found")
	}
}

func (s Server) InfoHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		respond.Auto(res, req, http.StatusOK, e.Responses[http.StatusOK])
	}
}
