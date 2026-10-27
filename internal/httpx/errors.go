package httpx

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Status  int
	Code    string
	Message string
	Details []FieldError
}

func (e *AppError) Error() string { return e.Message }

func BadRequest(msg string) *AppError {
	return &AppError{Status: 400, Code: "BAD_REQUEST", Message: msg}
}
func Unauthorized(msg string) *AppError {
	return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: msg}
}
func Forbidden() *AppError {
	return &AppError{Status: 403, Code: "FORBIDDEN", Message: "you do not have permission to perform this action"}
}
func NotFound(what string) *AppError {
	return &AppError{Status: 404, Code: "NOT_FOUND", Message: what + " not found"}
}
func Conflict(msg string) *AppError { return &AppError{Status: 409, Code: "CONFLICT", Message: msg} }
func ValidationError(details []FieldError) *AppError {
	return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Validation failed", Details: details}
}
func RateLimited(msg string) *AppError {
	return &AppError{Status: 429, Code: "RATE_LIMITED", Message: msg}
}
func InternalError() *AppError {
	return &AppError{Status: 500, Code: "INTERNAL_ERROR", Message: "internal server error"}
}
