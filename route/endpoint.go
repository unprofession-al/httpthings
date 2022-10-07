package route

import (
	"fmt"
	"net/http"
	"strconv"
)

type Endpoint struct {
	Parameters  []Parameter                     `json:"query_params" yaml:"query_params"`
	Name        string                          `json:"name" yaml:"name"`
	Description string                          `json:"description" yaml:"description"`
	RequestBody interface{}                     `json:"request" yaml:"request"`
	Responses   map[int]interface{}             `json:"response" yaml:"response"`
	HandlerFunc func(Endpoint) http.HandlerFunc `json:"-" yaml:"-"`
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

func (e Endpoint) GetParamAsString(name string, r *http.Request) (string, bool) {
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

func (e Endpoint) GetParamAsInt(name string, r *http.Request) (int, bool) {
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
