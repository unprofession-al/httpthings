package main

import (
	"net/http"

	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) routeConfig() r.RouteConfig {
	return r.RouteConfig{
		Routes: map[string]r.RouteConfig{
			"test": {
				Endpoints: map[string]r.Endpoint{
					http.MethodGet: {
						HandlerFunc: s.TestHandler,
					},
				},
			},
			"todos": {
				Endpoints: map[string]r.Endpoint{
					http.MethodGet: {
						HandlerFunc: s.ListTodosHandler,
						Responses: map[int]interface{}{
							http.StatusOK: []Todo{},
						},
					},
					http.MethodPost: {
						HandlerFunc: s.AddTodoHandler,
						RequestBody: TodoRequest{},
						Responses: map[int]interface{}{
							http.StatusOK: Todo{},
						},
					},
				},
				Routes: map[string]r.RouteConfig{
					"{name | Name of the Todo to filter}": {
						Endpoints: map[string]r.Endpoint{
							http.MethodGet: {
								HandlerFunc: s.ShowTodoHandler,
								Responses: map[int]interface{}{
									http.StatusOK: Todo{},
								},
							},
							http.MethodPut: {
								HandlerFunc: s.FinishTodoHandler,
								Responses: map[int]interface{}{
									http.StatusOK: Todo{},
								},
							},
						},
					},
				},
			},
		},
	}
}
