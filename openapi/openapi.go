package openapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/unprofession-al/httpthings/respond"
)

// Doc represents an [OpenAPI Document] according to the [OpenAPI Specification].
//
// [OpenAPI Document]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#oasDocument
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Doc struct {
	OpenAPI      string                `json:"openapi" yaml:"openapi"`
	Info         Info                  `json:"info" yaml:"info"`
	Servers      []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        Paths                 `json:"paths" yaml:"paths"`
	Components   Components            `json:"components,omitempty" yaml:"components,omitEmpty"`
	Tags         []Tag                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// HandleHTTP renders a Doc as YAML or JSON, based on the ending of the
// requst path.
func (doc *Doc) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".yaml") || strings.HasSuffix(r.URL.Path, ".yml") {
		respond.YAML(w, http.StatusOK, doc)
	} else {
		respond.JSON(w, http.StatusOK, doc)
	}
}

// MarshalJSON is a bit of a hack to ensure that refereces are set properly in
// the context of a OpenAPI Document. See [this pull request] for details.
//
// [this pull request]: https://github.com/invopop/jsonschema/pull/45]
func (doc *Doc) MarshalJSON() ([]byte, error) {
	type alias Doc
	raw, err := json.MarshalIndent((*alias)(doc), "", "    ")
	out := bytes.ReplaceAll(raw, []byte("#/$defs/"), []byte("#/components/schemas/"))
	return out, err
}

// Tag represents an [Tag Object] according to the [OpenAPI Specification].
//
// [Tag Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#tagObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Tag struct {
	Name         string                `json:"name" yaml:"name"`
	Description  string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocumentation represents an [External Documentation Object] according to the [OpenAPI Specification].
//
// [External Documentation Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#externalDocumentationObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type ExternalDocumentation struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// Info represents an [Info Object] according to the [OpenAPI Specification].
//
// [Info Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#infoObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Info struct {
	Title          string  `json:"title" yaml:"title"`
	Description    string  `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string  `json:"version" yaml:"version"`
}

// Contact represents a [Contact Object] according to the [OpenAPI Specification].
//
// [Contact Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#contactObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License represents a [License Object] according to the [OpenAPI Specification].
//
// [License Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#licenseObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Server represents a [Server Object] according to the [OpenAPI Specification].
//
// [Server Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#serverObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariables `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// Server represents a [Server Variables Object] according to the [OpenAPI Specification].
//
// [Server Variables Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#serverVariablesObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type ServerVariables struct {
	Enum        []string `json:"enum" yaml:"enum"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// Paths represents a [Paths Object] according to the [OpenAPI Specification].
//
// [Paths Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#pathsObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Paths map[string]PathItem

// PathItem represents a [Path Item Object] according to the [OpenAPI Specification].
//
// [Path Item Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#pathItemObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
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

// Operation represents an [Operation Object] according to the [OpenAPI Specification].
//
// [Operation Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#operationObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
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
	Security     []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// Parameter represents a [Parameter Object] according to the [OpenAPI Specification].
//
// [Parameter Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#parameterObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Parameter struct {
	Name            string `json:"name" yaml:"name"`
	In              string `json:"in" yaml:"in"`
	Description     string `json:"description" yaml:"description"`
	Required        bool   `json:"required" yaml:"required"`
	Deprecated      bool   `json:"deprecated" yaml:"deprecated"`
	AllowEmptyValue bool   `json:"allowEmptyValue" yaml:"allowEmptyValue"`
	Schema          Schema `json:"schema" yaml:"schema"`
}

// SecurityRequirement represents a [Security Requirement Object] according to the [OpenAPI Specification].
//
// [Security Requirement Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#securityRequirementObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type SecurityRequirement map[string][]string

// Schema represents a [Schema Object] according to the [OpenAPI Specification].
//
// [Schema Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#schemaObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Schema struct {
	Type  string  `json:"type,omitempty" yaml:"type,omitempty"`
	Ref   string  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Items *Schema `json:"items,omitempty" yaml:"items,omitempty"`
}

// Responses represents a [Responses Object] according to the [OpenAPI Specification].
//
// [Responses Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#responsesObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Responses map[string]Response

// Response represents a [Response Object] according to the [OpenAPI Specification].
//
// [Response Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#responseObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Response struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     Content `json:"content" yaml:"content"`
}

// Request represents a [Request Object] according to the [OpenAPI Specification].
//
// [Request Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#requestObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Request struct {
	Description string  `json:"description,omitempty" yaml:"description"`
	Content     Content `json:"content" yaml:"content"`
	Required    bool    `json:"required" yaml:"required"`
}

// Content represents the payload of a [Request Object] or a [Response Object] according to the [OpenAPI Specification].
//
// [Response Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#responseObject
// [Request Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#requestObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type Content map[string]schemaDef

type schemaDef struct {
	Schema Schema `json:"schema" yaml:"schema"`
}

// Components represents a [Components Object] according to the [OpenAPI Specification].
//
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
// [Components Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#componentsObject
type Components struct {
	Schemas         jsonschema.Definitions `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	SecuritySchemes SecuritySchemes        `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

// SecuritySchemes is a map of [Security Scheme Object] according to the [OpenAPI Specification].
//
// [Security Scheme Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#securitySchemeObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type SecuritySchemes map[string]SecurityScheme

// SecurityScheme represents a [Security Scheme Object] according to the [OpenAPI Specification].
//
// [Security Scheme Object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#securitySchemeObject
// [OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
type SecurityScheme struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
}
