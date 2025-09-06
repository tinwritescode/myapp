package common

// ErrorResponse represents a detailed error response
type ErrorResponse struct {
	BaseResponse
	Details string `json:"details,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	BaseResponse
	ValidationErrors []ValidationError `json:"validation_errors"`
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{
		BaseResponse: BaseResponse{
			Success: false,
			Error:   "Validation failed",
		},
		ValidationErrors: errors,
	}
}

// NewDetailedErrorResponse creates a detailed error response
func NewDetailedErrorResponse(message, details, code string) ErrorResponse {
	return ErrorResponse{
		BaseResponse: BaseResponse{
			Success: false,
			Error:   message,
		},
		Details: details,
		Code:    code,
	}
}
