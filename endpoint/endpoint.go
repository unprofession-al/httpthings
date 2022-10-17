package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Parameters    []Parameter
	Name          string
	Description   string
	RequestBody   interface{}
	Responses     map[int]interface{}
	ErrorResponse ErrorResponse
	Auth          *Auth
	Handler       func(w http.ResponseWriter, r *http.Request)
	Tags          []string
	Hidden        bool
}

type ErrorResponse interface {
	Respond(status int, details string, w http.ResponseWriter, r *http.Request)
}

func (e *Endpoint) RegisterError(status int, details string) http.HandlerFunc {
	e.Responses[status] = e.ErrorResponse
	return func(w http.ResponseWriter, r *http.Request) {
		e.ErrorResponse.Respond(status, details, w, r)
	}
}

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

func (e *Endpoint) GetParamAsInt(name string, r *http.Request) (int, bool) {
	param := Parameter{}
	for _, p := range e.Parameters {
		if p.Name != name {
			continue
		}
		param = p
		break
	}
	val, ok := param.First(r)
	if !ok {
		return 0, false
	}
	out, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return out, true
}

type Endpoints map[Caller]*Endpoint

type Caller struct {
	Path, Method string
}

func (c *Endpoints) Add(path, method string, e *Endpoint) error {
	err := checkHTTPVerbs(method)
	if err != nil {
		return err
	}
	path, params := extractPathParams(path)
	e.Parameters = append(e.Parameters, params...)
	(*c)[Caller{Path: path, Method: method}] = e
	return nil
}

func (c *Endpoints) PopulateRouter(router *mux.Router) {
	for caller, e := range *c {
		handler := e.Handler
		if e.Auth != nil && e.Auth.MiddlewareInjector != nil {
			handler = e.Auth.MiddlewareInjector(*e, handler)
		}
		path := c.preparePath(caller.Path)
		method := caller.Method
		if strings.HasSuffix(path, "*/") {
			path = strings.TrimSuffix(path, "*/")
			router.PathPrefix(path).HandlerFunc(handler).Methods(method)
		} else {
			router.Path(path).HandlerFunc(handler).Methods(method)
		}
	}
}

func (c *Endpoints) preparePath(path string) string {
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
