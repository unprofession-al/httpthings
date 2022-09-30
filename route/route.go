package route

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
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
	tidyPath, pathParams := extractPathParams(path)
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
	// By sorting the (sub)routes alphabetically we ensure that the order
	// it which the routes are attached to the mux.Router avoids the shadowing
	// of specific routes by "catch all" routes. For example:
	//
	// func main() {
	//     router := mux.NewRouter().StrictSlash(true)
	//     router.HandleFunc("/test/{var}", varHandler)
	//     router.HandleFunc("/test/a", aHandler)          // this route would never be executed
	//     http.ListenAndServe(":8888", router)
	// }
	//
	paths := make([]string, 0, len(c.Routes))
	for p := range c.Routes {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		subpath := fmt.Sprintf("%s%s/", path, strings.Trim(p, "/"))
		routes, err := expandRoutes(c.Routes[p], subpath)
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
		} else {
			router.Path(path).HandlerFunc(f(route.Endpoint)).Methods(route.Method)
		}
	}
}

func (r Routes) GetActionName(path, method string) (action string, found bool) {
	found = false
	for _, route := range r {
		if strings.ToLower(route.Method) != strings.ToLower(method) {
			continue
		}
		if !matchPath(route.Path, path) {
			continue
		}
		action = route.Endpoint.Name
		found = true
		break
	}
	return action, found
}

func matchPath(pattern, path string) bool {
	subExpr := regexp.MustCompile(`\{[a-zA-Z0-9-_]\}`)
	expression := subExpr.ReplaceAllString(pattern, `[a-zA-Z0-9-_%]*`)
	pathExpr, err := regexp.Compile(expression)
	if err != nil {
		return false
	}
	return pathExpr.MatchString(path)
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
