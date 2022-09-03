package respond

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Auto reads the 'accept' request header and tries to respond automatically with the appropriate
// 'content-type'. This currently works for 'text/yaml', everything else will be threaded as 'application/json'.
func Auto(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	switch req.Header.Get("Accept") {
	case "text/yaml":
		YAML(res, req, code, data)
	default:
		JSON(res, req, code, data)
	}

}

func YAML(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	res.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	out, err := yaml.Marshal(data)
	if err != nil {
		out = []byte("--- error: failed while rendering data to yaml")
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(out)
}

func JSON(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		out = []byte("{ 'error': 'failed while rendering data to json' }")
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(out)
}
