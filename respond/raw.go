package respond

import "net/http"

func Raw(res http.ResponseWriter, code int, data []byte) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
}
