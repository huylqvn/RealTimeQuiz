package errs

import "net/http"

var ErrInvalidRequest = errInvalidRequest{}

type errInvalidRequest struct{}

func (errInvalidRequest) Error() string {
	return "missing required fields or invalid request body"
}
func (errInvalidRequest) StatusCode() int {
	return http.StatusBadRequest
}

var ErrNotFound = errNotFound{}

type errNotFound struct{}

func (errNotFound) Error() string {
	return "record not exist"
}
func (errNotFound) StatusCode() int {
	return http.StatusBadRequest
}

var ErrSomethingWentWrong = errSomethingWentWrong{}

type errSomethingWentWrong struct{}

func (errSomethingWentWrong) Error() string {
	return "something went wrong"
}
func (errSomethingWentWrong) StatusCode() int {
	return http.StatusInternalServerError
}

var (
	ErrPermissionDenied = errPermissionDenied{}
	ErrAuthorization    = errAuthorization{}
	ErrNoContent        = errNoContent{}
)

type errPermissionDenied struct{}

func (e errPermissionDenied) Error() string {
	return "permission denied"
}
func (e errPermissionDenied) StatusCode() int {
	return http.StatusForbidden
}

type errAuthorization struct{}

func (e errAuthorization) Error() string {
	return "authorization failed"
}
func (e errAuthorization) StatusCode() int {
	return http.StatusUnauthorized
}

type errNoContent struct{}

func (e errNoContent) Error() string {
	return "no content"
}

func (e errNoContent) StatusCode() int {
	return http.StatusNoContent
}
