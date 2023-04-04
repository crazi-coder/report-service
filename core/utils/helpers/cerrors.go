package helpers

import "errors"

// ErrorCode represents the error code type
type ErrorCode int

// ErrUnAuthorized is used for returning custom error messages if the user is not authorized to perform the operation.
var ErrUnAuthorized = errors.New("user is not authorized to perform the operation")

// ErrPageLimitExceeded is used for returning custom error messages if the page limit is exceeded.
var ErrPageLimitExceeded = errors.New("maximum page limit exceeded")

// ErrRouteAlreadyLinked is used for returning custom error messages if the route already linked to the other Fe.
var ErrRouteAlreadyLinked = errors.New("route is already linked with other FieldExecutive")

const (

	// ErrCodeDataNotFound indicates the data is not found.
	ErrCodeDataNotFound ErrorCode = iota + 25000
	// ErrCodeStatusBadRequest indicates the request data send by client is not valid.
	ErrCodeStatusBadRequest
	// ErrCodeServerError indicates the server is not responding correctly.
	ErrCodeServerError
	//ErrPageLimitExceededError indicated that the requested page does not exist
	ErrPageLimitExceededError
	// ErrCodeUnauthorized indicates the user is not authorized to perform the operation.
	ErrCodeUnauthorized
	//ErrUUIDInvalid indicates the UUID is not valid
	ErrUUIDInvalid
	// ErrDateValidation indicates if the date validation is failed
	ErrDateValidation
)
