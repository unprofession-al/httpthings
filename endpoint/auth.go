package endpoint

import "net/http"

// Auth describes a [Security Schemes] according to the [OpenAPI Specification].
// It can be then linked to an [Endpoint]. If a MiddlewareInjector is provided,
// [Endpoint.PopulateRouter] will wrap the Handler of the endpoint in the middleware.
//
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
// [Security Schemes]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#componentsSecuritySchemes
type Auth struct {
	Name               string
	Type               string
	Scheme             string
	MiddlewareInjector AuthMiddlewareInjector
}

// AuthMiddlewareInjector is used to provide a middleware to the [Auth] struct.
// When populating the router, this function is called to wrap the handle func.
// Besides the [http.HandlerFunc] the AuthMiddlewareInjector function also takes
// the [Endpoint] itself as a parameter. This allows to make the middleware aware of
// the [Endpoint] it warps. This can be helpful if details of the endpoints
// are used in the auth process.
type AuthMiddlewareInjector func(Endpoint, http.HandlerFunc) http.HandlerFunc
