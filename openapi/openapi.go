package openapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/route"
)

type Spec struct {
	OpenAPI    string     `json:"openapi" yaml:"openapi"`
	Info       info       `json:"info" yaml:"info"`
	Servers    []server   `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths      paths      `json:"paths" yaml:"paths"`
	Components components `json:"components" yaml:"components"`
}

type info struct {
	Title          string  `json:"title" yaml:"title"`
	Description    string  `json:"description" yaml:"description"`
	TermsOfService string  `json:"termsOfService" yaml:"termsOfService"`
	Contact        contact `json:"contact" yaml:"contact"`
	License        license `json:"license" yaml:"license"`
	Version        string  `json:"version" yaml:"version"`
}

type contact struct {
	Name  string `json:"name" yaml:"name"`
	URL   string `json:"url" yaml:"url"`
	Email string `json:"email" yaml:"email"`
}

type license struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type server struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Variables serverVariablen
}

type paths map[string]endpoints

type endpoints map[string]endpoint

type endpoint struct {
	Description string              `json:"description,omitempty" yaml:"description"`
	Tags        []string            `json:"tags,omitempty" yaml:"tags"`
	RequestBody *request            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]response `json:"responses" yaml:"responses"`
	Parameters  []parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func newEndpoint(e route.Endpoint) (endpoint, []*jsonschema.Schema) {
	params := []parameter{}
	for _, p := range e.Parameters {
		param := parameter{
			Name:        p.Name,
			In:          p.Location.String(),
			Description: p.Description,
			Required:    true,
			Schema:      schema{Type: p.Content},
		}
		params = append(params, param)
	}
	body, bSchema := newRequest(e.RequestBody)
	responses, rSchemas := newResponses(e.Responses)
	out := endpoint{
		Description: e.Description,
		Responses:   responses,
		RequestBody: body,
		Parameters:  params,
	}
	schemas := append(rSchemas, bSchema)
	return out, schemas
}

type parameter struct {
	Name            string `json:"name" yaml:"name"`
	In              string `json:"in" yaml:"in"`
	Description     string `json:"description" yaml:"description"`
	Required        bool   `json:"required" yaml:"required"`
	Deprecated      bool   `json:"deprecated" yaml:"deprecated"`
	AllowEmptyValue bool   `json:"allowEmptyValue" yaml:"allowEmptyValue"`
	Schema          schema `json:"schema" yaml:"schema"`
}

type schema struct {
	Type  string  `json:"type,omitempty" yaml:"type,omitempty"`
	Ref   string  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Items *schema `json:"items,omitempty" yaml:"items,omitempty"`
}

func newSchema(t string, v interface{}) schema {
	array := reflect.TypeOf(v).Kind() == reflect.Slice
	out := &schema{}
	fill := out
	if array {
		out.Type = "array"
		out.Items = &schema{}
		fill = out.Items
	}
	switch t {
	case "bool":
		fill.Type = "boolean"
	case "string":
		fill.Type = t
	case "integer":
		fill.Type = t
	default:
		fill.Ref = t
	}

	return *out
}

type responses map[string]response

func newResponses(in map[int]interface{}) (responses, []*jsonschema.Schema) {
	out := responses{}
	schemas := []*jsonschema.Schema{}

	if len(in) == 0 {
		code := http.StatusOK
		resp, schema := newResponse(code, "")
		fmt.Println(schema)
		schemas = append(schemas, schema)
		out[fmt.Sprint(code)] = *resp
	}
	for code, data := range in {
		resp, schema := newResponse(code, data)
		schemas = append(schemas, schema)
		out[fmt.Sprint(code)] = *resp
	}
	return out, schemas
}

type response struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     content `json:"content" yaml:"content"`
}

func newResponse(code int, in interface{}) (*response, *jsonschema.Schema) {
	if in == nil {
		return nil, nil
	}
	r := jsonschema.Reflector{
		ReferenceRoot: "#/components/schemas/",
	}
	schema := r.Reflect(in)
	nameTokens := strings.SplitN(reflect.TypeOf(in).String(), ".", 2)
	var reference string
	if len(nameTokens) < 2 {
		reference = nameTokens[0]
	} else {
		reference = fmt.Sprintf("#/components/schemas/%s", nameTokens[1])
	}
	resp := &response{
		Description: statusText(code),
		Content: content{
			"application/json": {
				Schema: newSchema(reference, in),
			},
		},
	}
	return resp, schema
}

type request struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     content `json:"content" yaml:"content"`
	Required    bool    `json:"required" yaml:"required"`
}

func newRequest(in interface{}) (*request, *jsonschema.Schema) {
	if in == nil {
		return nil, nil
	}
	r := jsonschema.Reflector{
		ReferenceRoot: "#/components/schemas/",
	}
	schema := r.Reflect(in)
	nameTokens := strings.SplitN(reflect.TypeOf(in).String(), ".", 2)
	var reference string
	if len(nameTokens) < 2 {
		reference = nameTokens[0]
	} else {
		reference = fmt.Sprintf("#/components/schemas/%s", nameTokens[1])
	}
	req := &request{
		Content: content{
			"application/json": {
				Schema: newSchema(reference, in),
			},
		},
	}
	return req, schema
}

type content map[string]schemaDef

type schemaDef struct {
	Schema schema `json:"schema" yaml:"schema"`
}

type components struct {
	Schemas         jsonschema.Definitions `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	SecuritySchemes securitySchemes        `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

type securitySchemes map[string]securityScheme

type securityScheme struct{}

// https://apitools.dev/swagger-parser/online/
// https://www.thecodebuzz.com/swagger-openapi-3-0-sample-json-example-jwt-basic-auth/
