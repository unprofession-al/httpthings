package openapi

import (
	"net/http"
)

func (s *Spec) GetHandler() http.HandlerFunc {
	out, err := s.AsJSON()
	code := http.StatusOK
	if err != nil {
		out = []byte("{ 'error': 'failed while rendering data to json' }")
		code = http.StatusInternalServerError
	}
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.WriteHeader(code)
		res.Write(out)
	}
}
