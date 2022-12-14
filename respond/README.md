<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# respond

```go
import "github.com/unprofession-al/httpthings/respond"
```

Package respond provides functions to easily write http responses to the client

## Index

- [Constants](<#constants>)
- [func Auto(res http.ResponseWriter, req *http.Request, code int, data interface{}, headers ...map[string]string) error](<#func-auto>)
- [func JSON(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error](<#func-json>)
- [func Raw(res http.ResponseWriter, code int, data []byte, headers ...map[string]string)](<#func-raw>)
- [func YAML(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error](<#func-yaml>)


## Constants

```go
const (
    ContentTypeYAML = "text/yaml; charset=utf-8"        // default Content-Type when text/yaml is requested
    ContentTypeJSON = "application/json; charset=utf-8" // default Content-Type when json is requested
    ContentTypeRaw  = "text/plain; charset=utf-8"       // default Content-Type when rendering raw bytes
)
```

## func Auto

```go
func Auto(res http.ResponseWriter, req *http.Request, code int, data interface{}, headers ...map[string]string) error
```

Auto reads the 'accept' request header and tries to respond automatically with the appropriate 'content\-type'. This currently works for 'text/yaml', everything else will be threaded as 'application/json'.

## func JSON

```go
func JSON(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error
```

JSON uses the standard library to render the data provided as a JSON document, consult the \[docs\] to learn about on how to control the resulting output.

\[docs\]: https://pkg.go.dev/encoding/json

## func Raw

```go
func Raw(res http.ResponseWriter, code int, data []byte, headers ...map[string]string)
```

Raw writes plain bytes into the response and sets 'text/plain' as content type header if no "Content\-Type" header is provided.

## func YAML

```go
func YAML(res http.ResponseWriter, code int, data interface{}, headers ...map[string]string) error
```

YAML uses 'github.com/invopop/yaml' to render the data provided as a YAML document. Head to the \[official documentation\] to learn about the available tags to by used on the struct to control the output.

\[official documentation\]: https://github.com/invopop/yaml



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
