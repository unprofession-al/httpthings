package route

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Parameter struct {
	Name        string   `json:"name" yaml:"name"`
	Location    Location `json:"location" yaml:"location"`
	Required    bool     `json:"required" yaml:"required"`
	Default     *string  `json:"default" yaml:"default"`
	Description string   `json:"description" yaml:"description"`
	Content     string   `json:"content" yaml:"content"`
}

func (p Parameter) Get(r *http.Request) ([]string, bool) {
	if p.Location == LocationPath {
		v, ok := mux.Vars(r)[p.Name]
		if !ok && p.Default != nil {
			return []string{*p.Default}, true
		} else if !ok {
			return []string{}, false
		}
		return []string{v}, true
	} else if p.Location == LocationHeader {
		v := r.Header.Get(p.Name)
		if v == "" && p.Default != nil {
			return []string{*p.Default}, true
		} else if v == "" {
			return []string{}, false
		}
		return []string{v}, true
	} else if p.Location == LocationCookie {
		for _, cookie := range r.Cookies() {
			if cookie.Name == p.Name {
				return []string{cookie.Value}, true
			}
		}
		return []string{}, false
	} else if p.Location == LocationQuery {
		v, ok := r.URL.Query()[p.Name]
		if !ok && p.Default != nil {
			return []string{*p.Default}, true
		} else if !ok {
			return []string{}, false
		}
		return v, true
	}
	return []string{}, false
}

func (p Parameter) First(r *http.Request) (string, bool) {
	v, ok := p.Get(r)
	if len(v) == 0 {
		return "", ok
	}
	return v[0], ok
}

type Location int

const (
	LocationQuery Location = Location(iota)
	LocationHeader
	LocationPath
	LocationCookie
)

var locationText = map[Location]string{
	LocationQuery:  "query",
	LocationHeader: "header",
	LocationPath:   "path",
	LocationCookie: "cookie",
}

func NewLocation(in string) Location {
	in = strings.ToLower(in)
	for m, text := range locationText {
		if text == in {
			return Location(m)
		}
	}
	return LocationQuery
}

// String returns a string representation of the mode.
func (l Location) String() string {
	return locationText[l]
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
