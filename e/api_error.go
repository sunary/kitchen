package e

// APIError ...
type APIError struct {
	Code    Code // A standard grpc error code.
	Message string
	Err     error
}

// Error ...
func Error(code Code, message string, errs ...error) *APIError {
	if message == "" {
		message = code.String()
	}

	var err error
	if len(errs) > 0 {
		err = errs[0]
	}

	return &APIError{
		Code:    code,
		Err:     err,
		Message: message,
	}
}

// Error ...
func (e *APIError) Error() string {
	return e.Message
}
