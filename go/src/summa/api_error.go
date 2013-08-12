package summa

import (
	"net/http"
)

const (
	INTERNAL_ERROR = "Internal application error"
)

type apiError interface {
	Error() string
	Code() int
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
