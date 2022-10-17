package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// An Endpoint correlates pretty much with an [Operation Object] from the
// [OpenAPI Specification]. In addition to the meta data required by the
// [Operation Object] it also holds the [http.HandlerFunc] and provides some
// helpful mechanisms to work with HTTP errors as well as with the requests
// [Parameter]s.
//
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
// [Operation Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#operationObject
type Endpoint struct {
	// Name of the endpoint
	Name string
	// Description can be used to provide more detailes on the endpoint
	Description string
	// Parameters which can be provided to this endpoint via request
	Parameters []Parameter
	// RequestBody can be used to provide an example of the request body type
	// expected to be provided to this endpoint. When generating an OpenAPI Document
	// This value is then used to generate the schema to describe the request body.
	RequestBody interface{}
	// Responses are similar to the RequestBody field, but for resposes. The keys of
	// the map represent the HTTP status codes associated with the response.
	// Usually it is sufficient to describe the success responses here and then use
	// the RegisterError method for the bad exist cases.
	Responses map[int]interface{}
	// ErrorResponse can be used to provide a HTTP error 'template' in order to have
	// consistent errors across the API.
	ErrorResponse ErrorResponse
	// Auth references the auth method used for this endpoint.
	Auth *Auth
	// Handler is the actual [http.HandlerFunc] executed when this endpoint is called
	Handler http.HandlerFunc
	// Tags map directly to the tags field of the Operation Object according to the
	// OpenAPI Spec.
	Tags []string
	// Hidden prevents the endpoint from being represented in the OpenAPI document.
	Hidden bool
}

// ErrorResponse in an interface that can be implemented to ensure that all HTTP
// errors returned by an API are structured in the same manner. This helps to ensure
// that the usage of the API is as easy as possible. The ErrorResponse is used by
// [Endpoint.RegisterError] to build the [http.HandlerFunc].
type ErrorResponse interface {
	Respond(status int, details string, w http.ResponseWriter, r *http.Request)
}

// RegisterError takes a status code as well as some details on this error and
// registers a new [Response] to the [Endpoint]. As a result, a [http.HandlerFunc]
// is returned which can then be used to in the [Endpoint]s handler itself.
//
// Using these [http.HandlerFunc]s for Error handling in the handler function of the
// [Endpoint] is by no means neccessory but helps to ensure that all possible exit
// code are documented in the Responses map.
func (e *Endpoint) RegisterError(status int, details string) http.HandlerFunc {
	if e.ErrorResponse != nil {
		e.Responses[status] = e.ErrorResponse
		return func(w http.ResponseWriter, r *http.Request) {
			e.ErrorResponse.Respond(status, details, w, r)
		}
	} else {
		e.Responses[status] = details
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte(details))
		}
	}
}

// GetParamAsString fetches a the specified parameter from wherever it is stored in the
// given request as a string. If the value is found, it is returned
// as the first return value and `true` as the second return value. If the value cannot
// be found `false` will be returned as second return value.
func (e *Endpoint) GetParamAsString(name string, r *http.Request) (string, bool) {
	param := Parameter{}
	for _, p := range e.Parameters {
		if p.Name != name {
			continue
		}
		param = p
		break
	}
	return param.First(r)
}

// GetParamAsInt fetches a the specified parameter from wherever it is stored in the
// given request as an int. If the value is found, it is converted to an int and returned
// as the first return value and `true` as the second return value. If the value cannot
// be found or an error occures during conversion, `0` and `false` will be returned.
func (e *Endpoint) GetParamAsInt(name string, r *http.Request) (int, bool) {
	val, ok := e.GetParamAsString(name, r)
	if !ok {
		return 0, false
	}
	out, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return out, true
}

// Endpoints is a collection of references to an [Endpoint]. [Caller] is used to
// uniquely identify an [Endpoint]
type Endpoints map[Caller]*Endpoint

// A Caller specifies how a [http.Request] is mapped to an [Endpoint] by
// using the request path and HTTP method.
type Caller struct {
	Path, Method string
}

// Add appends an [Endpoint] to [Endpoints] at a given path and HTTP method.
func (c *Endpoints) Add(path, method string, e *Endpoint) error {
	err := checkHTTPVerbs(method)
	if err != nil {
		return err
	}
	path = preparePath(path)
	path, params := extractPathParams(path)
	caller := Caller{Path: path, Method: method}
	if _, exists := (*c)[caller]; exists {
		return fmt.Errorf("cannot add endpoint '%s', path '%s' with method '%s' already exists",
			e.Name, path, method)
	}
	e.Parameters = append(e.Parameters, params...)
	(*c)[caller] = e
	return nil
}

// PopulateRouter takes a reference to a [github.com/gorilla/mux.Router] and
// attaches all endpoints to it.
func (c *Endpoints) PopulateRouter(router *mux.Router) {
	for caller, e := range *c {
		handler := e.Handler
		if e.Auth != nil && e.Auth.MiddlewareInjector != nil {
			handler = e.Auth.MiddlewareInjector(*e, handler)
		}
		method, path := caller.Method, caller.Path
		if strings.HasSuffix(path, "*/") {
			path = strings.TrimSuffix(path, "*/")
			router.PathPrefix(path).HandlerFunc(handler).Methods(method)
		} else {
			router.Path(path).HandlerFunc(handler).Methods(method)
		}
	}
}

func preparePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	return path
}

func checkHTTPVerbs(verb string) error {
	allowed := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	legit := false
	for _, a := range allowed {
		if verb == a {
			legit = true
			break
		}
	}
	if !legit {
		return fmt.Errorf("'%s' is not an allowed HTTP verb", verb)
	}
	return nil
}
