/*
Package openapi provides a "good enough" implementaion of the [OpenAPI Specification].
It allows to generate a valid OpenAPI Document from [github.com/unprofession-al/httpthings/endpoint.Endpoints].

Most frameworks take the approach to generate code from documentation. This package
in conjunction with [github.com/unprofession-al/httpthings/endpoint] takes a different approach
and tries to generate the documentation from actual code while also tries to generate some
handy benefits from the additional code written. See [this discussion] for a bunch of oppinions
on the approaches to this topic.

[OpenAPI Specification]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md
[this discussion]: https://github.com/go-kit/kit/issues/185
*/
package openapi
