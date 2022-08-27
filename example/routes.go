package main

import (
	"github.com/invopop/jsonschema"
	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) routes() r.Route {
	return r.Route{
		Routes: map[string]r.Route{
			"test": {
				Handlers: map[string]*r.Endpoint{
					"GET": {
						HandlerFunc: s.TestHandler,
						Parameters:  []*r.Parameter{s.getParam("format")},
					},
				},
			},
			"openapi": {
				Handlers: map[string]*r.Endpoint{
					"GET": {
						HandlerFunc: s.OpenAPIHandler,
						Parameters:  []*r.Parameter{s.getParam("format")},
					},
				},
			},
			"todos": {
				Handlers: map[string]*r.Endpoint{
					"GET": {
						HandlerFunc: s.ListTodosHandler,
						Responses: map[string]*jsonschema.Schema{
							"200": jsonschema.Reflect([]Todo{}),
						},
						Parameters: []*r.Parameter{s.getParam("format")},
					},
					"POST": {
						HandlerFunc: s.AddTodoHandler,
						Parameters:  []*r.Parameter{s.getParam("format")},
						RequestBody: jsonschema.Reflect(TodoRequest{}),
						Responses: map[string]*jsonschema.Schema{
							"200": jsonschema.Reflect(Todo{}),
						},
					},
				},
				Routes: map[string]r.Route{
					"{name}": {
						Handlers: map[string]*r.Endpoint{
							"GET": {
								HandlerFunc: s.ShowTodoHandler,
								Parameters: []*r.Parameter{
									s.getParam("format"),
									s.getParam("name"),
								},
							},
							"PUT": {
								HandlerFunc: s.FinishTodoHandler,
								Parameters: []*r.Parameter{
									s.getParam("format"),
									s.getParam("name"),
								},
							},
						},
					},
				},
			},
		},
	}
}
