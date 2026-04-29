package errs

type ParseRequestError struct {
	Field   string
	Value   string
	Message string
}

func (e *ParseRequestError) Error() string {
	return e.Message
}

func NewParseRequestError(field, val, msg string) *ParseRequestError {
	return &ParseRequestError{
		Field:   field,
		Value:   val,
		Message: msg,
	}
}
