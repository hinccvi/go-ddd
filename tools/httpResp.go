package tools

import (
	"net/http"
)

type (
	message string

	response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	Success message = "success"
	Created message = "created"
	Updated message = "updated"
	Deleted message = "deleted"
	Error   message = "error"
)

func generateStatusCode(code int) int {
	if code > http.StatusNetworkAuthenticationRequired {
		code = http.StatusBadRequest
	}

	return code
}
