package route

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Handler struct {
	Parameters  []*Parameter            `json:"query_params" yaml:"query_params"`
	Desc        string                  `json:"description" yaml:"description"`
	RequestBody interface{}             `json:"request" yaml:"request"`
	Responses   map[string]interface{}  `json:"response" yaml:"response"`
	Func        func() http.HandlerFunc `json:"-" yaml:"-"`
}

func (h Handler) append(router *mux.Router, path, method string) (isCatchAll bool) {
	if strings.HasSuffix(path, "*/") {
		path = strings.TrimSuffix(path, "*/")

		f := h.Func
		if f == nil {
			f = notImplemented
		}
		router.PathPrefix(path).HandlerFunc(f()).Methods(method)
		return true
	}

	f := h.Func
	if f == nil {
		f = notImplemented
	}
	router.Path(path).HandlerFunc(f()).Methods(method)
	return false
}
