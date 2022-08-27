package route

import (
	"fmt"
	"strings"
)

func (r Route) getPaths(base string) map[string]string {
	out := map[string]string{}
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))
	//docs[base] = &openapi3.PathItem{}

	for path, route := range r.Routes {
		path = fmt.Sprintf("%s%s/", base, strings.Trim(path, "/"))
		sub := route.getPaths(path)
		for k, v := range sub {
			out[k] = v
		}

	}
	return out
}

func (r Route) Doc(base string) string {
	return ""
}
