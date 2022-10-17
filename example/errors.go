package main

import (
	"net/http"

	"github.com/unprofession-al/httpthings/respond"
)

type HTTPError struct {
	Code    int    `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}

func (e HTTPError) Respond(status int, details string, w http.ResponseWriter, r *http.Request) {
	e.Code = status
	e.Message = details
	respond.Auto(w, r, status, e)
}
