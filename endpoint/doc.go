/*
Package endpoint is a small layer on top of [github.com/gorilla/mux] which allows to
define HTTP endpoints containing the handler function as well as some meta data. In
conjunction with package `openapi` a set of endpoints can then be rendered as a JSON
or YAML representation of the endpoints OpenAPI specification.

To make best use of the meta data provided package endpoint also tries to provide
some functionality which attempts to make use of the meta data provided when
implementing the handler itself.
*/
package endpoint
