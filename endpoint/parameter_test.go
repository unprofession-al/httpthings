package endpoint

import (
	"net/http"
	"sort"
	"testing"
)

func TestParameterGet(t *testing.T) {
	cases := map[string]struct {
		request      *http.Request
		param        Parameter
		expectedVals []string
		expectedOK   bool
	}{
		"Header preset once with default": {
			request: requestWithHeaderParameters(nil, map[string][]string{"X-Test": {"foo"}}),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Header preset once with no default": {
			request: requestWithHeaderParameters(nil, map[string][]string{"X-Test": {"foo"}}),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Header preset more than once with no default": {
			request: requestWithHeaderParameters(nil, map[string][]string{"X-Test": {"foo", "bar"}}),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		"Header preset more than once with default": {
			request: requestWithHeaderParameters(nil, map[string][]string{"X-Test": {"foo", "bar"}}),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		"Header not set with no default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{},
			expectedOK:   false,
		},
		"Header not set with default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bla"},
			expectedOK:   true,
		},
		"Query set once with default": {
			request: requestWithQueryParameters(nil, map[string][]string{"test": {"foo"}}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Query set once default": {
			request: requestWithQueryParameters(nil, map[string][]string{"test": {"foo"}}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Query set more than once with no default": {
			request: requestWithQueryParameters(nil, map[string][]string{"test": {"foo", "bar"}}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		"Query set more than once with default": {
			request: requestWithQueryParameters(nil, map[string][]string{"test": {"foo", "bar"}}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		"Query not set with no default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Required: true,
			},
			expectedVals: []string{},
			expectedOK:   false,
		},
		"Query not set with default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationQuery,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bla"},
			expectedOK:   true,
		},
		"Cookie set once with default": {
			request: requestWithCookieParameters(nil, map[string]string{"test": "foo"}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationCookie,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Cookie set once with no default": {
			request: requestWithCookieParameters(nil, map[string]string{"test": "foo"}),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationCookie,
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		"Cookie not set with no default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationCookie,
				Required: true,
			},
			expectedVals: []string{},
			expectedOK:   false,
		},
		"Cookie not set with default": {
			request: requestWithNoParameters(),
			param: Parameter{
				Name:     "test",
				Location: ParameterLocationCookie,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bla"},
			expectedOK:   true,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			vals, ok := c.param.Get(c.request)
			if ok != c.expectedOK {
				t.Errorf("ok is not as expected, have %t, need %t", ok, c.expectedOK)
			}
			if !equal(vals, c.expectedVals) {
				t.Errorf("values not as expected, have %v, need %v", vals, c.expectedVals)
			}

		})
	}
}

func requestWithNoParameters() *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	return r
}

func requestWithCookieParameters(r *http.Request, p map[string]string) *http.Request {
	if r == nil {
		r = requestWithNoParameters()
	}
	for k, v := range p {
		c := &http.Cookie{Name: k, Value: v}
		r.AddCookie(c)
	}
	return r
}

func requestWithHeaderParameters(r *http.Request, p map[string][]string) *http.Request {
	if r == nil {
		r = requestWithNoParameters()
	}
	r.Header = p
	return r
}

func requestWithQueryParameters(r *http.Request, p map[string][]string) *http.Request {
	if r == nil {
		r = requestWithNoParameters()
	}
	q := r.URL.Query()
	for k, vals := range p {
		for _, v := range vals {
			q.Add(k, v)
		}
	}
	r.URL.RawQuery = q.Encode()
	return r
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/*
type Parameter struct {
	Name        string            `json:"name" yaml:"name"`
	Location    ParameterLocation `json:"location" yaml:"location"`
	Required    bool              `json:"required" yaml:"required"`
	Default     string            `json:"default" yaml:"default"`
	Description string            `json:"description" yaml:"description"`
	Type        string            `json:"content" yaml:"content"`
}
*/
