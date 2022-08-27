package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type Endpoint struct {
	Parameters  []*Parameter                    `json:"query_params" yaml:"query_params"`
	Desc        string                          `json:"description" yaml:"description"`
	RequestBody *jsonschema.Schema              `json:"request" yaml:"request"`
	Responses   map[string]*jsonschema.Schema   `json:"response" yaml:"response"`
	HandlerFunc func(Endpoint) http.HandlerFunc `json:"-" yaml:"-"`
}

func (e Endpoint) append(router *mux.Router, path, method string) (isCatchAll bool) {
	if strings.HasSuffix(path, "*/") {
		path = strings.TrimSuffix(path, "*/")

		f := e.HandlerFunc
		if f == nil {
			f = notImplemented
		}
		router.PathPrefix(path).HandlerFunc(f(e)).Methods(method)
		return true
	}

	f := e.HandlerFunc
	if f == nil {
		f = notImplemented
	}
	router.Path(path).HandlerFunc(f(e)).Methods(method)
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

func (e Endpoint) ValidateRequestBody(document []byte) (*gojsonschema.Result, error) {
	schema, err := e.RequestBody.MarshalJSON()
	if err != nil {
		return nil, err
	}
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(document)
	return gojsonschema.Validate(schemaLoader, documentLoader)
}
