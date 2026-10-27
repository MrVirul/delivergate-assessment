package httpx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]FieldError, 0, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			msg := fmt.Sprintf("failed on '%s' tag", fe.Tag())
			// Map common tags
			switch fe.Tag() {
			case "required":
				msg = "is required"
			case "min":
				msg = "is too short or too small"
			case "max":
				msg = "is too long or too large"
			case "gte":
				msg = "must be >= " + fe.Param()
			case "lte":
				msg = "must be <= " + fe.Param()
			case "email":
				msg = "must be a valid email"
			case "oneof":
				msg = "must be one of: " + fe.Param()
			}
			details = append(details, FieldError{Field: field, Message: msg})
		}
		return ValidationError(details)
	}
	return BadRequest(err.Error())
}
