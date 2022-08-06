package respond

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Auto reads the 'accept' request header and tries to respond automatically with the appropriate
// 'content-type'. This currently works for 'text/yaml', everything else will be threaded as 'application/json'.
func Auto(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	var err error
	var errMesg []byte
	var out []byte

	f := req.Header.Get("Accept")
	if f == "text/yaml" {
		res.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		out, err = yaml.Marshal(data)
		errMesg = []byte("--- error: failed while rendering data to yaml")
	} else {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		out, err = json.MarshalIndent(data, "", "    ")
		errMesg = []byte("{ 'error': 'failed while rendering data to json' }")
	}

	if err != nil {
		out = errMesg
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(out)
}
