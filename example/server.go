package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
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
	fmt.Printf("Serving at http://%s\nPress CTRL-c to stop...\n", s.listener)
	log.Fatal(http.ListenAndServe(s.listener, s.handler))
}

func (s Server) respond(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	if code != http.StatusOK {
		fmt.Println(data)
	}
	var err error
	var errMesg []byte
	var out []byte

	f := s.param("format").First(req)
	if f == "yaml" {
		res.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		out, err = yaml.Marshal(data)
		errMesg = []byte("--- error: failed while rendering data to yaml")
	} else {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		out, err = json.Marshal(data)
		errMesg = []byte("{ 'error': 'failed while rendering data to json' }")
	}

	if err != nil {
		out = errMesg
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(out)
}

func (s Server) raw(res http.ResponseWriter, code int, data []byte) {
	res.WriteHeader(code)
	res.Write(data)
}
