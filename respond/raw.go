package respond

import "net/http"

// Raw writes plain bytes into the response and sets 'text/plain' as content type
// header.
func Raw(res http.ResponseWriter, code int, data []byte) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
}
