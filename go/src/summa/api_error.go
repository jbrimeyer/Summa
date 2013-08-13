package summa

import (
	"net/http"
)

const (
	INVALID_OR_MISSING = "Invalid or missing field"
	METHOD_NOT_ALLOWED = "Invalid request method"
	INTERNAL_ERROR     = "Internal application error"
)

type apiError interface {
	Error() string
	Code() int
	Data() apiResponseData
}

// 400
type badRequestError struct {
	s string
}

func (e *badRequestError) Error() string {
	return e.s
}

func (e *badRequestError) Code() int {
	return http.StatusBadRequest
}

func (e *badRequestError) Data() apiResponseData {
	return nil
}

// 401
type unauthorizedError struct {
	s string
}

func (e *unauthorizedError) Error() string {
	return e.s
}

func (e *unauthorizedError) Code() int {
	return http.StatusUnauthorized
}

func (e *unauthorizedError) Data() apiResponseData {
	return nil
}

// 403
type forbiddenError struct {
	s string
}

func (e *forbiddenError) Error() string {
	return e.s
}

func (e *forbiddenError) Code() int {
	return http.StatusForbidden
}

func (e *forbiddenError) Data() apiResponseData {
	return nil
}

// 404
type notFoundError struct {
	s string
}

func (e *notFoundError) Error() string {
	return e.s
}

func (e *notFoundError) Code() int {
	return http.StatusNotFound
}

func (e *notFoundError) Data() apiResponseData {
	return nil
}

// 405
type methodNotAllowedError struct {
}

func (e *methodNotAllowedError) Error() string {
	return METHOD_NOT_ALLOWED
}

func (e *methodNotAllowedError) Code() int {
	return http.StatusMethodNotAllowed
}

func (e *methodNotAllowedError) Data() apiResponseData {
	return nil
}

// 409
type conflictError struct {
	data apiResponseData
}

func (e *conflictError) Error() string {
	return INVALID_OR_MISSING
}

func (e *conflictError) Code() int {
	return http.StatusConflict
}

func (e *conflictError) Data() apiResponseData {
	return e.data
}

// 500
type internalServerError struct {
	s   string
	err error
}

func (e *internalServerError) Error() string {
	return INTERNAL_ERROR
}

func (e *internalServerError) Code() int {
	return http.StatusInternalServerError
}

func (e *internalServerError) Data() apiResponseData {
	return nil
}
