package endpoint

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

// Parameter represents the [Parameter Object] of the [OpenAPI Specification]. It also cames
// with some handy functions to extrat the parameters from a [http.Request].
//
// [Parameter Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#parameterObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#parameterObject
type Parameter struct {
	Name        string            `json:"name" yaml:"name"`
	Location    ParameterLocation `json:"location" yaml:"location"`
	Required    bool              `json:"required" yaml:"required"`
	Default     string            `json:"default" yaml:"default"`
	Description string            `json:"description" yaml:"description"`
	Type        string            `json:"content" yaml:"content"`
}

// Get returns values if a given [http.Request]. To do this it respects the
// Location field of the parameter. If the parameter in not present in the request,
// the Defalut will be returned.
func (p Parameter) Get(r *http.Request) ([]string, bool) {
	if p.Location == ParameterLocationPath {
		v, ok := mux.Vars(r)[p.Name]
		if !ok && p.Default != "" {
			return []string{p.Default}, true
		} else if !ok {
			return []string{}, false
		}
		return []string{v}, true
	} else if p.Location == ParameterLocationHeader {
		v := r.Header.Get(p.Name)
		if v == "" && p.Default != "" {
			return []string{p.Default}, true
		} else if v == "" {
			return []string{}, false
		}
		return []string{v}, true
	} else if p.Location == ParameterLocationCookie {
		for _, cookie := range r.Cookies() {
			if cookie.Name == p.Name {
				return []string{cookie.Value}, true
			}
		}
		return []string{}, false
	} else if p.Location == ParameterLocationQuery {
		v, ok := r.URL.Query()[p.Name]
		if !ok && p.Default != "" {
			return []string{p.Default}, true
		} else if !ok {
			return []string{}, false
		}
		return v, true
	}
	return []string{}, false
}

// First does the same as Get but only returns the first value
// found.
func (p Parameter) First(r *http.Request) (string, bool) {
	v, ok := p.Get(r)
	if len(v) == 0 {
		return "", ok
	}
	return v[0], ok
}

// ParameterLocation is used to describe where in a request a certain parameter
// can be found.
type ParameterLocation int

const (
	ParameterLocationQuery  ParameterLocation = ParameterLocation(iota) // parameter is located in the query string
	ParameterLocationHeader                                             // parameter is located in the request header
	ParameterLocationPath                                               // parameter is part of the request path
	ParameterLocationCookie                                             // parameter is stored in a cookie
)

var parameterLocationText = map[ParameterLocation]string{
	ParameterLocationQuery:  "query",
	ParameterLocationHeader: "header",
	ParameterLocationPath:   "path",
	ParameterLocationCookie: "cookie",
}

// String returns a string representation of the mode.
func (l ParameterLocation) String() string {
	return parameterLocationText[l]
}

func extractPathParams(path string) (tidy string, params []Parameter) {
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
			Location:    ParameterLocationPath,
			Required:    true,
			Default:     "",
			Description: desc,
			Type:        "string",
		}
		params = append(params, p)
	}
	return
}
