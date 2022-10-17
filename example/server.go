package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/unprofession-al/httpthings/endpoint"
	"github.com/unprofession-al/httpthings/openapi"
	"github.com/unprofession-al/httpthings/run"
)

type Server struct {
	listener string
	handler  http.Handler
	todos    TodoSet
	spec     openapi.Doc
	auth     *endpoint.Auth
}

func NewServer(listener, static string) (Server, error) {
	basicAuth := &endpoint.Auth{
		Name:               "BasicAuth",
		Type:               "http",
		Scheme:             "basic",
		MiddlewareInjector: WrapBasicAuth,
	}
	s := Server{
		listener: listener,
		todos:    *(NewTodoSet().Prepopulate()),
		auth:     basicAuth,
	}

	endpoints := &endpoint.Endpoints{}
	endpoints.Add("/api/v1/todos/", http.MethodGet, s.ListTodoEndpoint())
	endpoints.Add("/api/v1/todos/", http.MethodPost, s.AddTodoEndpoint())
	endpoints.Add("/api/v1/todos/{name | Name of the todo}/", http.MethodGet, s.ShowTodoEndpoint())
	endpoints.Add("/api/v1/todos/{name | Name of the todo}/", http.MethodPut, s.FinishTodoEndpoint())

	r := mux.NewRouter()
	endpoints.PopulateRouter(r)
	s.spec = openapi.FromEndpoints(*endpoints)
	s.spec.OpenAPI = "3.0.3"
	s.spec.Info.Version = "v1"
	s.spec.Info.Title = "Todo API"
	s.spec.Servers = []openapi.Server{{URL: fmt.Sprintf("http://%s", listener)}}
	r.Path("/openapi.json").HandlerFunc(s.spec.HandleHTTP)
	r.Path("/openapi.yaml").HandlerFunc(s.spec.HandleHTTP)
	s.handler = alice.New(cors.Default().Handler).Then(r)
	return s, nil
}

func (s *Server) run() {
	err := run.Run(run.DetectRunMode(), s.listener, s.handler, func(log string) { fmt.Printf("INFO: %s", log) })
	if err != nil {
		fmt.Println(err)
	}
}
