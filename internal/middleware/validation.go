package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tinwritescode/myapp/internal/dto/common"
)

// ValidationTag represents the type of validation tag
type ValidationTag string

const (
	TagRequired           ValidationTag = "required"
	TagRequiredIf         ValidationTag = "required_if"
	TagRequiredUnless     ValidationTag = "required_unless"
	TagRequiredWith       ValidationTag = "required_with"
	TagRequiredWithAll    ValidationTag = "required_with_all"
	TagRequiredWithout    ValidationTag = "required_without"
	TagRequiredWithoutAll ValidationTag = "required_without_all"
	TagEmail              ValidationTag = "email"
	TagMin                ValidationTag = "min"
	TagMax                ValidationTag = "max"
	TagLen                ValidationTag = "len"
	TagNumeric            ValidationTag = "numeric"
	TagAlpha              ValidationTag = "alpha"
	TagAlphanum           ValidationTag = "alphanum"
	TagURL                ValidationTag = "url"
	TagUUID               ValidationTag = "uuid"
	TagOneOf              ValidationTag = "oneof"
	TagGTE                ValidationTag = "gte"
	TagLTE                ValidationTag = "lte"
	TagGT                 ValidationTag = "gt"
	TagLT                 ValidationTag = "lt"
	TagNe                 ValidationTag = "ne"
	TagEq                 ValidationTag = "eq"
	TagContains           ValidationTag = "contains"
	TagExcludes           ValidationTag = "excludes"
	TagStartswith         ValidationTag = "startswith"
	TagEndswith           ValidationTag = "endswith"
	TagIP                 ValidationTag = "ip"
	TagIPv4               ValidationTag = "ipv4"
	TagIPv6               ValidationTag = "ipv6"
	TagMAC                ValidationTag = "mac"
	TagHostname           ValidationTag = "hostname"
	TagFQDN               ValidationTag = "fqdn"
	TagUnique             ValidationTag = "unique"
	TagJSON               ValidationTag = "json"
	TagJWT                ValidationTag = "jwt"
	TagLowercase          ValidationTag = "lowercase"
	TagUppercase          ValidationTag = "uppercase"
	TagDatetime           ValidationTag = "datetime"
	TagTimezone           ValidationTag = "timezone"
	TagBoolean            ValidationTag = "boolean"
)

// formatValidationErrors formats validator errors to be more readable like Zod
func formatValidationErrors(err error) []common.ValidationError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []common.ValidationError
		for _, e := range validationErrors {
			fieldName := strings.ToLower(e.Field())
			message := getValidationMessage(e.Tag(), fieldName, e.Param())
			errors = append(errors, common.ValidationError{
				Field:   fieldName,
				Message: message,
				Value:   fmt.Sprintf("%v", e.Value()),
			})
		}
		return errors
	}
	// Fallback for non-validation errors
	return []common.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func getValidationMessage(tag, field, param string) string {
	switch ValidationTag(tag) {
	case TagRequired, TagRequiredIf, TagRequiredUnless, TagRequiredWith, TagRequiredWithAll, TagRequiredWithout, TagRequiredWithoutAll:
		return fmt.Sprintf("%s is required", field)
	case TagEmail:
		return fmt.Sprintf("%s must be a valid email address", field)
	case TagMin:
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case TagMax:
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case TagLen:
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case TagNumeric:
		return fmt.Sprintf("%s must be a number", field)
	case TagAlpha:
		return fmt.Sprintf("%s must contain only letters", field)
	case TagAlphanum:
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case TagURL:
		return fmt.Sprintf("%s must be a valid URL", field)
	case TagUUID:
		return fmt.Sprintf("%s must be a valid UUID", field)
	case TagOneOf:
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case TagGTE:
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case TagLTE:
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case TagGT:
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case TagLT:
		return fmt.Sprintf("%s must be less than %s", field, param)
	case TagNe:
		return fmt.Sprintf("%s must not be equal to %s", field, param)
	case TagEq:
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case TagContains:
		return fmt.Sprintf("%s must contain %s", field, param)
	case TagExcludes:
		return fmt.Sprintf("%s must not contain %s", field, param)
	case TagStartswith:
		return fmt.Sprintf("%s must start with %s", field, param)
	case TagEndswith:
		return fmt.Sprintf("%s must end with %s", field, param)
	case TagIP:
		return fmt.Sprintf("%s must be a valid IP address", field)
	case TagIPv4:
		return fmt.Sprintf("%s must be a valid IPv4 address", field)
	case TagIPv6:
		return fmt.Sprintf("%s must be a valid IPv6 address", field)
	case TagMAC:
		return fmt.Sprintf("%s must be a valid MAC address", field)
	case TagHostname:
		return fmt.Sprintf("%s must be a valid hostname", field)
	case TagFQDN:
		return fmt.Sprintf("%s must be a valid FQDN", field)
	case TagUnique:
		return fmt.Sprintf("%s must be unique", field)
	case TagJSON:
		return fmt.Sprintf("%s must be valid JSON", field)
	case TagJWT:
		return fmt.Sprintf("%s must be a valid JWT", field)
	case TagLowercase:
		return fmt.Sprintf("%s must be lowercase", field)
	case TagUppercase:
		return fmt.Sprintf("%s must be uppercase", field)
	case TagDatetime:
		return fmt.Sprintf("%s must be a valid datetime", field)
	case TagTimezone:
		return fmt.Sprintf("%s must be a valid timezone", field)
	case TagBoolean:
		return fmt.Sprintf("%s must be a boolean value", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// BindJSON validates and binds JSON request body to the provided struct
// Returns true if binding was successful, false if there was an error
func BindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		validationErrors := formatValidationErrors(err)
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(validationErrors))
		return false
	}
	return true
}

// BindQuery validates and binds query parameters to the provided struct
// Returns true if binding was successful, false if there was an error
func BindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse("Invalid query parameters"))
		return false
	}
	return true
}
