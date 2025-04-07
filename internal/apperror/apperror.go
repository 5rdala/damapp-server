package apperror

import "net/http"

type HTTPStatusCode = int

const (
	ErrCodeDataNotFound      HTTPStatusCode = http.StatusNotFound
	ErrCodeInvalidCredential HTTPStatusCode = http.StatusUnauthorized
	ErrCodeInternal          HTTPStatusCode = http.StatusInternalServerError
	ErrCodeBadRequest        HTTPStatusCode = http.StatusBadRequest
)

type AppError struct {
	Code    HTTPStatusCode
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code HTTPStatusCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
