package endpoint

import (
	"net/http"
	"sort"
	"testing"
)

func TestParameterGet(t *testing.T) {
	cases := []struct {
		name         string
		request      *http.Request
		param        Parameter
		expectedVals []string
		expectedOK   bool
	}{
		{
			name:    "Header preset once with no default",
			request: &http.Request{Header: map[string][]string{"X-Test": {"foo"}}},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		{
			name:    "Header preset once with no default",
			request: &http.Request{Header: map[string][]string{"X-Test": {"foo"}}},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{"foo"},
			expectedOK:   true,
		},
		{
			name:    "Header preset more than once with no default",
			request: &http.Request{Header: map[string][]string{"X-Test": {"foo", "bar"}}},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		{
			name:    "Header preset more than once with default",
			request: &http.Request{Header: map[string][]string{"X-Test": {"foo", "bar"}}},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bar", "foo"},
			expectedOK:   true,
		},
		{
			name:    "Header not preset with no default",
			request: &http.Request{},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Required: true,
			},
			expectedVals: []string{},
			expectedOK:   false,
		},
		{
			name:    "Header not preset with default",
			request: &http.Request{},
			param: Parameter{
				Name:     "X-Test",
				Location: ParameterLocationHeader,
				Default:  "bla",
				Required: true,
			},
			expectedVals: []string{"bla"},
			expectedOK:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
