package main

import (
	"net/http"

	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) params() map[string]*r.Parameter {
	formatDefault := "json"
	return map[string]*r.Parameter{
		"format": {
			Name:     "Accept",
			Default:  &formatDefault,
			Location: r.LocationHeader,
			Desc:     "format of the output, can be 'yaml' or 'json'",
		},
		"name": {
			Name:     "name",
			Location: r.LocationPath,
			Default:  nil,
			Desc:     "name of the todo",
		},
	}
}

func (s Server) getParam(name string) *r.Parameter {
	return s.params()[name]
}

func (s Server) extractParam(name string, r *http.Request) (string, bool) {
	param, ok := s.params()[name]
	if !ok {
		return "", false
	}
	return param.First(r)
}
