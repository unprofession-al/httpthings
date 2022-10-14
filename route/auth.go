package route

import "net/http"

type Auth struct {
	Name              string
	Type              string
	Scheme            string
	MiddlewareWrapper func(Endpoint, http.HandlerFunc) http.HandlerFunc
}
