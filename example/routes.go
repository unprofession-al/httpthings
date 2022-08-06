package main

import (
	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) routes() r.Route {
	return r.Route{
		R: r.Routes{
			"openapi.html": {
				H: r.Handlers{"GET": r.Handler{F: s.OpenAPIHandler, Q: []*r.QueryParam{s.param("format")}}},
			},
			"todos": {
				H: r.Handlers{
					"GET": r.Handler{F: s.ListTodosHandler, Q: []*r.QueryParam{s.param("format")}},
					"POST": r.Handler{
						F:   s.AddTodoHandler,
						Q:   []*r.QueryParam{s.param("format")},
						Req: Todo{},
						Res: Todo{},
					},
				},
				R: r.Routes{
					"{todo}": {
						H: r.Handlers{
							"GET": r.Handler{F: s.ShowTodoHandler, Q: []*r.QueryParam{s.param("format")}},
							"PUT": r.Handler{F: s.FinishTodoHandler, Q: []*r.QueryParam{s.param("format")}},
						},
					},
				},
			},
		},
	}
}
