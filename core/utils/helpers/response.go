package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Error   string    `json:"error,omitempty"`
}

type Response interface {
	Error(code ErrorCode, message string, err error) errorResponse
}

type response struct{}

func NewResponse() Response {

	return &response{}
}

func (r *response) Error(code ErrorCode, message string, e error) errorResponse {
	// Few messages are fix, which should not be changed.
	switch code {
	case ErrCodeServerError:
		message = fmt.Sprintf("Has encountered a situation it doesn't know how to handle. %s", message)
	case ErrCodeStatusBadRequest:
		message = "The server will not process the request due to something that is perceived to be an error."
	}
	if gin.Mode() == gin.DebugMode {
		if e != nil {
			return errorResponse{Code: code, Message: message, Error: e.Error()}
		}
		return errorResponse{Code: code, Message: message, Error: ""}
	}
	return errorResponse{Code: code, Message: message}
}
