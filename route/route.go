package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type RouteConfig struct {
	Endpoints map[string]Endpoint
	Routes    map[string]RouteConfig
}

type Routes []Route

func NewRoutes(c RouteConfig, base string) (Routes, error) {
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))
	return expandRoutes(c, base)
}

func expandRoutes(c RouteConfig, path string) (Routes, error) {
	r := []Route{}
	tidyPath, pathParams := tidyPath(path)
	path = tidyPath
	err := checkHTTPVerbs(c.Endpoints)
	if err != nil {
		return r, err
	}
	for m, ep := range c.Endpoints {
		ep.Parameters = append(ep.Parameters, pathParams...)
		route := Route{
			Method:   strings.ToUpper(m),
			Path:     path,
			Endpoint: ep,
		}
		r = append(r, route)
	}
	for p, route := range c.Routes {
		subpath := fmt.Sprintf("%s%s/", path, strings.Trim(p, "/"))
		routes, err := expandRoutes(route, subpath)
		if err != nil {
			return r, err
		}
		r = append(r, routes...)
	}
	return r, err
}

type Route struct {
	Method   string
	Path     string
	Endpoint Endpoint
}

func (r Routes) PopulateRouter(router *mux.Router) {
	for _, route := range r {
		path := route.Path
		f := route.Endpoint.HandlerFunc
		if f == nil {
			f = notImplemented
		}
		if strings.HasSuffix(path, "*/") {
			path = strings.TrimSuffix(path, "*/")
			router.PathPrefix(path).HandlerFunc(f(route.Endpoint)).Methods(route.Method)
		}
		router.Path(path).HandlerFunc(f(route.Endpoint)).Methods(route.Method)
	}
}

func notImplemented(e Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotImplemented)
		out := "Not Yet Implemented\n"
		res.Write([]byte(out))
	}
}

func checkHTTPVerbs(h map[string]Endpoint) error {
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
	for verb := range h {
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
	}
	return nil
}
