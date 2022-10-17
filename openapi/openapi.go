package openapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/respond"
)

type OpenAPI struct {
	OpenAPI      string                `json:"openapi" yaml:"openapi"`
	Info         Info                  `json:"info" yaml:"info"`
	Servers      []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        Paths                 `json:"paths" yaml:"paths"`
	Components   Components            `json:"components,omitempty" yaml:"components,omitEmpty"`
	Tags         []Tag                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

func (oapi *OpenAPI) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".yaml") || strings.HasSuffix(r.URL.Path, ".yml") {
		respond.YAML(w, http.StatusOK, oapi)
	} else {
		respond.JSON(w, http.StatusOK, oapi)
	}
}

func (oapi *OpenAPI) MarshalJSON() ([]byte, error) {
	type alias OpenAPI
	raw, err := json.MarshalIndent((*alias)(oapi), "", "    ")
	out := bytes.ReplaceAll(raw, []byte("#/$defs/"), []byte("#/components/schemas/"))
	return out, err
}

type Tag struct {
	Name         string                `json:"name" yaml:"name"`
	Description  string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type ExternalDocumentation struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

type Info struct {
	Title          string  `json:"title" yaml:"title"`
	Description    string  `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string  `json:"version" yaml:"version"`
}

type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariables `json:"variables,omitempty" yaml:"variables,omitempty"`
}

type ServerVariables struct {
	Enum        []string `json:"enum" yaml:"enum"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

type Paths map[string]PathItem

type PathItem struct {
	Ref         string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation  `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []Server    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type Operation struct {
	Tags         []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   []Parameter           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *Request              `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    map[string]Response   `json:"responses" yaml:"responses"`
	Deprecated   bool                  `json:"deprecated" yaml:"deprecated"`
	Security     SecurityRequirements  `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
}

type Parameter struct {
	Name            string `json:"name" yaml:"name"`
	In              string `json:"in" yaml:"in"`
	Description     string `json:"description" yaml:"description"`
	Required        bool   `json:"required" yaml:"required"`
	Deprecated      bool   `json:"deprecated" yaml:"deprecated"`
	AllowEmptyValue bool   `json:"allowEmptyValue" yaml:"allowEmptyValue"`
	Schema          Schema `json:"schema" yaml:"schema"`
}

type SecurityRequirements []SecurityRequirement

type SecurityRequirement map[string][]string

type Schema struct {
	Type  string  `json:"type,omitempty" yaml:"type,omitempty"`
	Ref   string  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Items *Schema `json:"items,omitempty" yaml:"items,omitempty"`
}

type Responses map[string]Response

type Response struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     Content `json:"content" yaml:"content"`
}

type Request struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     Content `json:"content" yaml:"content"`
	Required    bool    `json:"required" yaml:"required"`
}

type Content map[string]SchemaDef

type SchemaDef struct {
	Schema Schema `json:"schema" yaml:"schema"`
}

type Components struct {
	Schemas         jsonschema.Definitions `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	SecuritySchemes SecuritySchemes        `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

type SecuritySchemes map[string]SecurityScheme

type SecurityScheme struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
}

// https://apitools.dev/swagger-parser/online/
// https://www.thecodebuzz.com/swagger-openapi-3-0-sample-json-example-jwt-basic-auth/
