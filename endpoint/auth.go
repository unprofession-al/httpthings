package endpoint

import "net/http"

type AuthMiddlewareInjector func(Endpoint, http.HandlerFunc) http.HandlerFunc

type Auth struct {
	Name               string
	Type               string
	Scheme             string
	MiddlewareInjector AuthMiddlewareInjector
}
