package openapi

import (
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/route"
)

type Config struct {
	Title             string
	Version           string
	Description       string
	TermsOfService    string
	ContactName       string
	ContactURL        string
	ContactEmail      string
	LicenseName       string
	LicenseURL        string
	ServerURL         string
	ServerDescription string
}

func New(c Config, r route.Routes, base string) Spec {
	spec := Spec{
		OpenAPI: "3.0.3",
		Info: info{
			Title:          c.Title,
			Version:        c.Version,
			Description:    c.Description,
			TermsOfService: c.TermsOfService,
			Contact: contact{
				Name:  c.ContactName,
				URL:   c.ContactURL,
				Email: c.ContactEmail,
			},
			License: license{
				Name: c.LicenseName,
				URL:  c.LicenseURL,
			},
		},
		Paths: paths{},
		Servers: []server{
			{
				URL:         c.ServerURL,
				Description: c.ServerDescription,
			},
		},
	}
	schemas := []*jsonschema.Schema{}
	secSchemes := map[string]securityScheme{}
	for _, route := range r {
		if _, ok := spec.Paths[route.Path]; !ok {
			spec.Paths[route.Path] = endpoints{}
		}
		var epSchemas []*jsonschema.Schema
		var secScheme map[string]securityScheme
		spec.Paths[route.Path][strings.ToLower(route.Method)], epSchemas, secScheme = newEndpoint(route.Endpoint)
		schemas = append(schemas, epSchemas...)
		for k, v := range secScheme {
			secSchemes[k] = v
			fmt.Printf("name: %s, val: %v\n", k, v)
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
