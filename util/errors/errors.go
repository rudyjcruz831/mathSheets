package errors

import (
	"fmt"
	"net/http"
	"time"
)

// Type holds a type string and integer code for the error
type Type string

// "Set" of valid errorTypes
const (
	Authorization        Type = "AUTHORIZATION"          // Authentication Failures -
	BadRequest           Type = "BAD_REQUEST"            // Validation errors / BadInput
	Conflict             Type = "CONFLICT"               // Already exists (eg, create account with existent email) - 409
	Internal             Type = "INTERNAL"               // Server (500) and fallback errors
	NotFound             Type = "NOT_FOUND"              // For not finding resource
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"      // for uploading tons of JSON, or an image over the limit - 413
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"    // For long running handlers
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
)

type MathSheetsError struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func UnauthorizedError(message string) *MathSheetsError {
	return &MathSheetsError{
		Timestamp: time.Now(),
		Message:   message,
		Status:    http.StatusUnauthorized,
		Error:     string(Authorization),
	}
}

func NewBadRequestError(message string) *MathSheetsError {
	return &MathSheetsError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   string(BadRequest),
	}
}

func NewInternalServerError(message string) *MathSheetsError {
	return &MathSheetsError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   string(Internal),
	}
}

func NewNotFound(name, value string) *MathSheetsError {
	return &MathSheetsError{
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
		Status:  http.StatusNotFound,
		Error:   string(NotFound),
	}
}

func NewUnsupportedMediaType(message string) *MathSheetsError {
	return &MathSheetsError{
		Message: message,
		Status:  http.StatusUnsupportedMediaType,
		Error:   string(UnsupportedMediaType),
	}
}

// New Conflict to create an error for 409
func NewConflict(name, value string) *MathSheetsError {
	return &MathSheetsError{
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
		Status:  http.StatusConflict,
		Error:   string(Conflict),
	}
}

func NewServiceUnavailable() *MathSheetsError {
	return &MathSheetsError{
		Message: "Service unavailable or time out",
		Status:  http.StatusServiceUnavailable,
		Error:   string(ServiceUnavailable),
	}
}

func NewAuthorization(message string) *MathSheetsError {
	return &MathSheetsError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   string(Authorization),
	}
}

// {
// 	"message":"Did not find User --- Some error down call chain",
// 	"status":404,
// 	"error":"NOT_FOUND"
// }
