package main

import (
	r "github.com/unprofession-al/httpthings/route"
)

func (s Server) params() map[string]r.Parameter {
	return map[string]r.Parameter{}
}

func (s Server) getParam(name string) r.Parameter {
	return s.params()[name]
}
