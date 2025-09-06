package common

import "fmt"

// ERROR_CODE represents error codes used throughout the application
type ERROR_CODE int

const (
	EMAIL_ALREADY_USED ERROR_CODE = iota
	INVALID_CREDENTIALS
	USER_NOT_FOUND
	INVALID_TOKEN
	TOKEN_EXPIRED
	UNAUTHORIZED
	FORBIDDEN
	VALIDATION_ERROR
	INTERNAL_SERVER_ERROR
	BAD_REQUEST
	NOT_FOUND
	CONFLICT
	TOO_MANY_REQUESTS
)

// String returns the string representation of the error code
func (e ERROR_CODE) String() string {
	switch e {
	case EMAIL_ALREADY_USED:
		return "EMAIL_ALREADY_USED"
	case INVALID_CREDENTIALS:
		return "INVALID_CREDENTIALS"
	case USER_NOT_FOUND:
		return "USER_NOT_FOUND"
	case INVALID_TOKEN:
		return "INVALID_TOKEN"
	case TOKEN_EXPIRED:
		return "TOKEN_EXPIRED"
	case UNAUTHORIZED:
		return "UNAUTHORIZED"
	case FORBIDDEN:
		return "FORBIDDEN"
	case VALIDATION_ERROR:
		return "VALIDATION_ERROR"
	case INTERNAL_SERVER_ERROR:
		return "INTERNAL_SERVER_ERROR"
	case BAD_REQUEST:
		return "BAD_REQUEST"
	case NOT_FOUND:
		return "NOT_FOUND"
	case CONFLICT:
		return "CONFLICT"
	case TOO_MANY_REQUESTS:
		return "TOO_MANY_REQUESTS"
	default:
		return "UNKNOWN_ERROR"
	}
}

// ErrorResponse represents a detailed error response
type ErrorResponse struct {
	BaseResponse
	Code string `json:"code,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	ErrorResponse
	ValidationErrors []ValidationError `json:"validation_errors"`
}

// AppError represents an application error with error code
type AppError struct {
	Code    ERROR_CODE
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code ERROR_CODE, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewErrorResponseWithCode creates an error response with error code
func NewErrorResponseWithCode(code ERROR_CODE, message string) ErrorResponse {
	return ErrorResponse{
		BaseResponse: BaseResponse{
			Success: false,
			Error:   message,
		},
		Code: code.String(),
	}
}

func NewValidationErrorResponse(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{
		ErrorResponse: ErrorResponse{
			BaseResponse: BaseResponse{
				Success: false,
				Error:   "Validation failed",
			},
			Code: VALIDATION_ERROR.String(),
		},
		ValidationErrors: errors,
	}
}
