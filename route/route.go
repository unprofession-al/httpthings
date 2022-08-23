package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Handlers map[string]*Endpoint
	Routes   map[string]Route
}

func (r Route) Populate(router *mux.Router, base string) error {
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))

	err := checkHTTPVerbs(r.Handlers)
	if err != nil {
		return err
	}

	for verb, handler := range r.Handlers {
		handler.append(router, base, strings.ToUpper(verb))
	}

	for path, route := range r.Routes {
		path = fmt.Sprintf("%s%s/", base, strings.Trim(path, "/"))
		err = route.Populate(router, path)
	}
	return err
}

func notImplemented(e Endpoint) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotImplemented)
		out := "Not Yet Implemented\n"
		res.Write([]byte(out))
	}
}

func checkHTTPVerbs(h map[string]*Endpoint) error {
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
