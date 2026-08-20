package models

import (
	"time"

	"github.com/google/uuid"
)

// ResponseEnvelope represents the standard KIVU API JSON envelope.
type ResponseEnvelope struct {
	Data  interface{} `json:"data"`
	Error *APIError   `json:"error"`
}

// APIError details API errors consistently.
type APIError struct {
	Message string    `json:"message"`
	Code    int       `json:"code,omitempty"`
	Details string    `json:"details,omitempty"`
	Time    time.Time `json:"time"`
}

// NewSuccessResponse creates a standard success envelope.
func NewSuccessResponse(data interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		Data:  data,
		Error: nil,
	}
}

// NewErrorResponse creates a standard error envelope.
func NewErrorResponse(message string, code int) ResponseEnvelope {
	return ResponseEnvelope{
		Data: nil,
		Error: &APIError{
			Message: message,
			Code:    code,
			Time:    time.Now().UTC(),
		},
	}
}

// MustParseUUID parses a UUID string or returns Nil.
func MustParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}
