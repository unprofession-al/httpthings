package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

type Route struct {
	Connect *Handler
	Delete  *Handler
	Get     *Handler
	Head    *Handler
	Options *Handler
	Patch   *Handler
	Post    *Handler
	Put     *Handler
	Trace   *Handler
	Routes  map[string]Route
}

func (r Route) asPathItem(base string) *openapi3.PathItem {
	out := &openapi3.PathItem{}
	return out
}

func (r Route) Populate(router *mux.Router, base string) {
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))

	if r.Connect != nil {
		r.Connect.append(router, base, "CONNECT")
	}
	if r.Delete != nil {
		r.Delete.append(router, base, "DELETE")
	}
	if r.Get != nil {
		r.Get.append(router, base, "GET")
	}
	if r.Head != nil {
		r.Head.append(router, base, "HEAD")
	}
	if r.Options != nil {
		r.Options.append(router, base, "OPTIONS")
	}
	if r.Patch != nil {
		r.Patch.append(router, base, "PATCH")
	}
	if r.Post != nil {
		r.Post.append(router, base, "POST")
	}
	if r.Put != nil {
		r.Put.append(router, base, "PUT")
	}
	if r.Trace != nil {
		r.Trace.append(router, base, "TRACE")
	}

	for path, route := range r.Routes {
		path = fmt.Sprintf("%s%s/", base, strings.Trim(path, "/"))
		route.Populate(router, path)
	}
}

func notImplemented(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotImplemented)
	out := "Not Yet Implemented\n"
	res.Write([]byte(out))
}
