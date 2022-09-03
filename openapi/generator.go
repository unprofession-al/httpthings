package openapi

import (
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/route"
)

type Config struct {
	Title          string
	Version        string
	Description    string
	TermsOfService string
	ContactName    string
	ContactURL     string
	ContactEmail   string
	LicenseName    string
	LicenseURL     string
}

type generator struct {
	spec    *Spec
	schemas jsonschema.Definitions
}

func New(c Config, r route.Routes, base string) *Spec {
	g := &generator{
		spec: &Spec{
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
		},
		schemas: jsonschema.Definitions{},
	}
	for _, route := range r {
		if _, ok := g.spec.Paths[route.Path]; !ok {
			g.spec.Paths[route.Path] = endpoints{}
		}
		g.spec.Paths[route.Path][strings.ToLower(route.Method)] = g.newEndpoint(route.Endpoint)
	}
	g.spec.Components.Schemas = g.schemas
	return g.spec
}
