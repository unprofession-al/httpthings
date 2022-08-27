package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/respond"
	"github.com/unprofession-al/httpthings/route"
)

func (s Server) OpenAPIHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		out := s.routes().Doc(s.base)
		respond.Auto(res, req, http.StatusOK, out)
	}
}

func (s Server) ListTodosHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		respond.Auto(res, req, http.StatusOK, s.todos.AsSlice())
	}
}

func (s Server) TestHandler(e route.Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		schema := jsonschema.Reflect(Todo{})
		out, err := schema.MarshalJSON()
		if err != nil {
			respond.Auto(res, req, http.StatusInternalServerError, err)
			return
		}
		respond.Raw(res, http.StatusOK, out)
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
			respond.Auto(res, req, http.StatusNotAcceptable, fmt.Sprintf("todo not provided"))
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
		valRes, err := e.ValidateRequestBody(b)
		if err != nil {
			respond.Auto(res, req, http.StatusInternalServerError, "could not validate Request body")
			return
		}
		if !valRes.Valid() {
			errDesc := "The request body is not valid. Errors are: "
			for _, desc := range valRes.Errors() {
				errDesc += desc.String() + "... "
			}
			respond.Auto(res, req, http.StatusNotAcceptable, errDesc)
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
			respond.Auto(res, req, http.StatusNotAcceptable, fmt.Sprintf("todo not provided"))
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
