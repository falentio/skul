package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger = log.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

type Response interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ResponseError interface {
	Response

	Error() string
}

type HttpResponse[T any] struct {
	Code    int            `json:"code"`
	Data    T              `json:"data"`
	Cookies []*http.Cookie `json:"-"`
}

func (o *HttpResponse[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, c := range o.Cookies {
		http.SetCookie(w, c)
	}
	w.WriteHeader(o.Code)
	if o.Code == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		HandleError(w, r, err)
	}
}

func NewResponse[T any](data T, code int) *HttpResponse[T] {
	return &HttpResponse[T]{code, data, nil}
}

func NewOK[T any](data T) *HttpResponse[T] {
	return NewResponse(data, http.StatusOK)
}

func NewCreated[T any](data T) *HttpResponse[T] {
	return NewResponse(data, http.StatusCreated)
}

func NewNoContent() *HttpResponse[any] {
	return NewResponse[any](nil, http.StatusNoContent)
}

type HttpResponsePaginate[T any] struct {
	*HttpResponse[T]
	Data []T `json:"data"`
	Page Page
}

type Page struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

func NewPaginate[T any](datas []T, page Page) HttpResponsePaginate[T] {
	return HttpResponsePaginate[T]{
		HttpResponse: &HttpResponse[T]{
			Code: http.StatusOK,
		},
		Data: datas,
		Page: page,
	}
}

type HttpError struct {
	Errors  map[string]string `json:"errors,omitempty"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if re, ok := err.(ResponseError); !ok {
		logger.Error().Err(err).Msg("internal server error")
		HandleError(w, r, NewInternalServerError(nil, "internal server error"))
	} else {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		re.ServeHTTP(w, r)
	}
}

func (he *HttpError) Error() string {
	return he.Message
}

func (he *HttpError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, max-age=0, no-cache, no-store")
	w.WriteHeader(he.Code)
	if err := json.NewEncoder(w).Encode(he); err != nil {
		logger.Error().Err(err).Msg("failed to serve http")
		http.Error(w, "failed to marshall response body", http.StatusUnprocessableEntity)
	}
}

func NewError(errors map[string]string, message string, code int) ResponseError {
	return &HttpError{errors, message, code}
}

func NewBadRequest(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusBadRequest)
}

func NewUnauthorized(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusUnauthorized)
}

func NewForbidden(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusForbidden)
}

func NewConflict(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusConflict)
}

func NewUnprocessableEntity(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusUnprocessableEntity)
}

func NewNotFound(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusNotFound)
}

func NewInternalServerError(errors map[string]string, message string, a ...any) ResponseError {
	message = fmt.Sprintf(message, a...)
	return NewError(errors, message, http.StatusInternalServerError)
}
