// Package respond provides functions to easily write http responses to the client
package respond

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/yaml"
)

const (
	ContentTypeYAML = "text/yaml; charset=utf-8"        // default Content-Type when text/yaml is requested
	ContentTypeJSON = "application/json; charset=utf-8" // default Content-Type when json is requested
	ContentTypeRaw  = "text/plain; charset=utf-8"       // default Content-Type when rendering raw bytes
)

// Auto reads the 'accept' request header and tries to respond automatically with the appropriate
// 'content-type'. This currently works for 'text/yaml', everything else will be threaded as
// 'application/json'.
func Auto(res http.ResponseWriter, req *http.Request, code int, data interface{}, headers ...map[string]string) error {
	switch req.Header.Get("Accept") {
	case "text/yaml":
		return YAML(res, code, data, headers...)
	default:
		return JSON(res, code, data, headers...)
	}
}

// YAML uses 'github.com/invopop/yaml' to render the data provided as a YAML document. Head to the
// [official documentation] to learn about the available tags to by used on the struct to control the
// output.
//
// [official documentation]: https://github.com/invopop/yaml
func YAML(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error {
	for k, v := range getHeaders(ContentTypeRaw, headers...) {
		res.Header().Add(k, v)
	}
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to yaml: %w", err)
	}
	res.WriteHeader(code)
	_, err = res.Write(out)
	return err
}

// JSON uses the standard library to render the data provided as a JSON document, consult the [docs]
// to learn about on how to control the resulting output.
//
// [docs]: https://pkg.go.dev/encoding/json
func JSON(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error {
	for k, v := range getHeaders(ContentTypeRaw, headers...) {
		res.Header().Add(k, v)
	}
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal to json: %w", err)
	}
	res.WriteHeader(code)
	_, err = res.Write(out)
	return err
}

// Raw writes plain bytes into the response and sets 'text/plain' as content type
// header if no "Content-Type" header is provided.
func Raw(res http.ResponseWriter, code int, data []byte, headers ...map[string]string) {
	for k, v := range getHeaders(ContentTypeRaw, headers...) {
		res.Header().Add(k, v)
	}
	res.WriteHeader(code)
	res.Write(data)
}

func getHeaders(defaultContentType string, headers ...map[string]string) map[string]string {
	out := map[string]string{}
	hasContentType := false
	for _, h := range headers {
		for k, v := range h {
			out[k] = v
			if k == "Content-Type" {
				hasContentType = true
			}
		}
	}
	if !hasContentType {
		out["Content-Type"] = defaultContentType
	}
	return out
}
