// Package respond provides functions to easily write http responses to the client
package respond

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
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

// YAML uses 'gopkg.in/yaml.v3' to render the data provided as a YAML document. Head to the
// [official documentation] to learn about the available tags to by used on the struct to controll the
// output.
//
// [official documentation]: https://pkg.go.dev/gopkg.in/yaml.v3
func YAML(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error {
	setHeaders(res, ContentTypeYAML, headers...)
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to marshal to yaml: %w", err)
	}
	res.WriteHeader(code)
	res.Write(out)
	return nil
}

// JSON uses the standard library to render the data provided as a JSON document, consult the [docs]
// to learn about on how to controll the resulting output.
//
// [docs]: https://pkg.go.dev/encoding/json
func JSON(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error {
	setHeaders(res, ContentTypeJSON, headers...)
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to marshal to json: %w", err)
	}
	res.WriteHeader(code)
	res.Write(out)
	return nil
}

// Raw writes plain bytes into the response and sets 'text/plain' as content type
// header if no "Content-Type" header is provided.
func Raw(res http.ResponseWriter, code int, data []byte, headers ...map[string]string) {
	setHeaders(res, ContentTypeRaw, headers...)
	res.WriteHeader(code)
	res.Write(data)
}

func setHeaders(res http.ResponseWriter, defaultContentType string, headers ...map[string]string) {
	hasContentType := false
	for _, h := range headers {
		for k, v := range h {
			res.Header().Set(k, v)
			if k == "Content-Type" {
				hasContentType = true
			}
		}
	}
	if !hasContentType {
		res.Header().Set("Content-Type", defaultContentType)
	}
}
