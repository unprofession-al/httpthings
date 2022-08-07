package main

import (
	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) routes() r.Route {
	return r.Route{
		Routes: map[string]r.Route{
			"openapi": {
				Get: &r.Handler{
					Func:       s.OpenAPIHandler,
					Parameters: []*r.Parameter{s.getParam("format")},
				},
			},
			"todos": {
				Get: &r.Handler{
					Func: s.ListTodosHandler,
					Responses: map[string]interface{}{
						"200": []Todo{},
					},
					Parameters: []*r.Parameter{s.getParam("format")},
				},
				Post: &r.Handler{
					Func:        s.AddTodoHandler,
					Parameters:  []*r.Parameter{s.getParam("format")},
					RequestBody: Todo{},
					Responses: map[string]interface{}{
						"200": Todo{},
					},
				},
				Routes: map[string]r.Route{
					"{name}": {
						Get: &r.Handler{
							Func:       s.ShowTodoHandler,
							Parameters: []*r.Parameter{s.getParam("format")},
						},
						Put: &r.Handler{
							Func:       s.FinishTodoHandler,
							Parameters: []*r.Parameter{s.getParam("format")},
						},
					},
				},
			},
		},
	}
}
