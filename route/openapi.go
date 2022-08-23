package route

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func (r Route) getPathItems(base string) openapi3.Paths {
	docs := openapi3.Paths{}
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))
	docs[base] = &openapi3.PathItem{}

	for path, route := range r.Routes {
		path = fmt.Sprintf("%s%s/", base, strings.Trim(path, "/"))
		sub := route.getPathItems(path)
		for k, v := range sub {
			docs[k] = v
		}

	}
	return docs
}

func (r Route) Doc(base string) openapi3.T {
	doc := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "ToDo API",
			Description: "REST APIs used for interacting with the ToDo Service",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/MarioCarrion/todo-api-microservice-example",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://0.0.0.0:9234",
			},
		},
		Paths: r.getPathItems(base),
	}

	return doc
}
