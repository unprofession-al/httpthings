package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unprofession-al/httpthings/openapi"
	"github.com/unprofession-al/httpthings/route"
	"github.com/unprofession-al/httpthings/run"
)

type Server struct {
	listener string
	base     string
	handler  *mux.Router
	todos    TodoSet
	spec     openapi.Spec
}

func NewServer(listener, static string) (Server, error) {
	ts := NewTodoSet()
	// Populate some initial tasks
	ts.Add(&Todo{
		Name:        "Task1",
		Description: "The First Task",
	})
	ts.Add(&Todo{
		Name:        "Task2",
		Description: "The Second Task",
	})

	s := Server{
		listener: listener,
		todos:    *ts,
		base:     "api/v1",
	}

	r := mux.NewRouter()

	routes, err := route.NewRoutes(s.routeConfig(), s.base)
	if err != nil {
		return s, err
	}
	routes.PopulateRouter(r)

	r.Path("/openapi.json").HandlerFunc(s.OpenAPIHandler)

	if static != "" {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(static)))
	}

	s.GenerateOpenAPI(routes)

	s.handler = r
	return s, nil
}

func (s *Server) run() {
	run.Run(run.DetectRunMode(), s.listener, s.handler, true)
}

func (s *Server) GenerateOpenAPI(r route.Routes) {
	c := openapi.Config{
		Title:          "Todo API Service",
		Version:        "v1",
		Description:    "Todo API service is an example project of github.com/unprofession-al/httpthings",
		TermsOfService: "https://very.unprofession.al/termsOfService",
		ContactName:    "unprofession.al",
		ContactURL:     "https://very.unprofession.al/",
		ContactEmail:   "noreply@unprofession.al",
		LicenseName:    "The MIT License",
		LicenseURL:     "https://mit-license.org/",
		ServerURL:      fmt.Sprintf("http://%s", s.listener),
	}
	s.spec = openapi.New(c, r, s.base)
}
