package openapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/endpoint"
)

func FromEndpoints(col endpoint.Endpoints) OpenAPI {
	spec := OpenAPI{
		Paths: Paths{},
	}
	schemas := []*jsonschema.Schema{}
	secSchemes := map[string]SecurityScheme{}
	for caller, endpoint := range col {
		if endpoint.Hidden {
			continue
		}
		path, ok := spec.Paths[caller.Path]
		if !ok {
			path = PathItem{}
		}
		o, epSchemas, secScheme := newOperation(endpoint, endpoint.Tags...)
		switch strings.ToUpper(caller.Method) {
		case http.MethodGet:
			path.Get = o
		case http.MethodPut:
			path.Put = o
		case http.MethodPost:
			path.Post = o
		case http.MethodDelete:
			path.Delete = o
		case http.MethodOptions:
			path.Options = o
		case http.MethodHead:
			path.Head = o
		case http.MethodPatch:
			path.Patch = o
		case http.MethodTrace:
			path.Trace = o
		}
		spec.Paths[caller.Path] = path
		schemas = append(schemas, epSchemas...)
		for k, v := range secScheme {
			secSchemes[k] = v
		}
	}
	defs := jsonschema.Definitions{}
	for _, schema := range schemas {
		if schema == nil || schema.Definitions == nil {
			continue
		}
		for key, def := range schema.Definitions {
			defs[key] = def
		}
	}
	spec.Components.Schemas = defs
	spec.Components.SecuritySchemes = secSchemes
	return spec
}

func newOperation(e *endpoint.Endpoint, tags ...string) (*Operation, []*jsonschema.Schema, SecuritySchemes) {
	params := []Parameter{}
	for _, p := range e.Parameters {
		param := Parameter{
			Name:        p.Name,
			In:          p.Location.String(),
			Description: p.Description,
			Required:    true,
			Schema:      Schema{Type: p.Type},
		}
		params = append(params, param)
	}
	body, bSchema := newRequest(e.RequestBody)
	responses, rSchemas := newResponses(e.Responses)
	out := &Operation{
		Summary:     e.Name,
		Description: e.Description,
		OperationID: e.Description,
		Responses:   responses,
		RequestBody: body,
		Parameters:  params,
		Tags:        tags,
	}
	schemas := append(rSchemas, bSchema)
	sec := map[string]SecurityScheme{}
	if e.Auth != nil {
		sec[e.Auth.Name] = SecurityScheme{Type: e.Auth.Type, Scheme: e.Auth.Scheme}
		out.Security = []SecurityRequirement{
			{e.Auth.Name: []string{}},
		}
	}
	return out, schemas, sec
}

func newSchema(t string, v interface{}) Schema {
	array := reflect.TypeOf(v).Kind() == reflect.Slice
	out := &Schema{}
	fill := out
	if array {
		out.Type = "array"
		out.Items = &Schema{}
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

func newResponses(in map[int]interface{}) (Responses, []*jsonschema.Schema) {
	out := Responses{}
	schemas := []*jsonschema.Schema{}

	if len(in) == 0 {
		code := http.StatusOK
		resp, schema := newResponse(code, "")
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

func newResponse(code int, in interface{}) (*Response, *jsonschema.Schema) {
	if in == nil {
		return nil, nil
	}
	schema := jsonschema.Reflect(in)
	nameTokens := strings.SplitN(reflect.TypeOf(in).String(), ".", 2)
	var reference string
	if len(nameTokens) < 2 {
		reference = nameTokens[0]
	} else {
		reference = fmt.Sprintf("#/components/schemas/%s", nameTokens[1])
	}
	resp := &Response{
		Description: statusText(code),
		Content: Content{
			"application/json": {
				Schema: newSchema(reference, in),
			},
		},
	}
	return resp, schema
}

func newRequest(in interface{}) (*Request, *jsonschema.Schema) {
	if in == nil {
		return nil, nil
	}
	schema := jsonschema.Reflect(in)
	nameTokens := strings.SplitN(reflect.TypeOf(in).String(), ".", 2)
	var reference string
	if len(nameTokens) < 2 {
		reference = nameTokens[0]
	} else {
		reference = fmt.Sprintf("#/components/schemas/%s", nameTokens[1])
	}
	req := &Request{
		Content: Content{
			"application/json": {
				Schema: newSchema(reference, in),
			},
		},
	}
	return req, schema
}
