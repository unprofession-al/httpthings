package main

import (
	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) routes() r.Route {
	return r.Route{
		Routes: map[string]r.Route{
			"openapi": {
				Handlers: map[string]*r.Handler{
					"GET": {
						Func:       s.OpenAPIHandler,
						Parameters: []*r.Parameter{s.getParam("format")},
					},
				},
			},
			"todos": {
				Handlers: map[string]*r.Handler{
					"GET": {
						Func: s.ListTodosHandler,
						Responses: map[string]interface{}{
							"200": []Todo{},
						},
						Parameters: []*r.Parameter{s.getParam("format")},
					},
					"POST": {
						Func:        s.AddTodoHandler,
						Parameters:  []*r.Parameter{s.getParam("format")},
						RequestBody: Todo{},
						Responses: map[string]interface{}{
							"200": Todo{},
						},
					},
				},
				Routes: map[string]r.Route{
					"{name}": {
						Handlers: map[string]*r.Handler{
							"GET": {
								Func:       s.ShowTodoHandler,
								Parameters: []*r.Parameter{s.getParam("format")},
							},
							"PUT": {
								Func:       s.FinishTodoHandler,
								Parameters: []*r.Parameter{s.getParam("format")},
							},
						},
					},
				},
			},
		},
	}
}
