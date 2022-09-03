package route

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Parameters  []Parameter                     `json:"query_params" yaml:"query_params"`
	Name        string                          `json:"name" yaml:"name"`
	Desciption  string                          `json:"description" yaml:"description"`
	RequestBody interface{}                     `json:"request" yaml:"request"`
	Responses   map[int]interface{}             `json:"response" yaml:"response"`
	HandlerFunc func(Endpoint) http.HandlerFunc `json:"-" yaml:"-"`
}

func (e *Endpoint) append(router *mux.Router, path, method string) (isCatchAll bool) {
	method = strings.ToUpper(method)
	_, pathParams := tidyPath(path)
	e.Parameters = append(e.Parameters, pathParams...)
	if strings.HasSuffix(path, "*/") {
		path = strings.TrimSuffix(path, "*/")
		f := e.HandlerFunc
		if f == nil {
			f = notImplemented
		}
		router.PathPrefix(path).HandlerFunc(f(*e)).Methods(method)
		return true
	}
	f := e.HandlerFunc
	if f == nil {
		f = notImplemented
	}
	router.Path(path).HandlerFunc(f(*e)).Methods(method)
	return false
}

func (e Endpoint) GetParams(r *http.Request) (map[string][]string, []error) {
	out := map[string][]string{}
	errs := []error{}

	ok := false
	for _, p := range e.Parameters {
		out[p.Name], ok = p.Get(r)
		if !ok && p.Required {
			errs = append(errs, fmt.Errorf("parameter '%s' is required but was not provided", p.Name))
		}
	}

	return out, errs
}

func tidyPath(path string) (tidy string, params []Parameter) {
	tidy = path
	rex := regexp.MustCompile(`\{(?P<param>.*)\}`)
	matches := rex.FindAllStringSubmatch(path, -1)
	for _, param := range matches {
		if len(param) < 2 {
			continue
		}
		whole := param[0]
		pair := strings.SplitN(param[1], "|", 2)
		key := strings.TrimSpace(pair[0])
		desc := key
		if len(pair) == 2 {
			desc = strings.TrimSpace(pair[1])
		}
		tidy = strings.ReplaceAll(tidy, whole, fmt.Sprintf("{%s}", key))
		p := Parameter{
			Name:        key,
			Location:    LocationPath,
			Required:    true,
			Default:     nil,
			Description: desc,
			Content:     "string",
		}
		params = append(params, p)
	}
	return
}
