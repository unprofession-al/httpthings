package main

import r "github.com/unprofession-al/httpthings/route"

func (s Server) params() map[string]*r.QueryParam {
	return map[string]*r.QueryParam{
		"format": {
			N:    "f",
			D:    "json",
			Desc: "format of the output, can be 'yaml' or 'json'",
		},
		"name": {
			N:    "n",
			D:    "",
			Desc: "name of the todo",
		},
	}
}

func (s Server) param(name string) *r.QueryParam {
	return s.params()[name]
}
