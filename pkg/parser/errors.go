package parser

import (
	"fmt"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// ParseError represents an error during annotation parsing
type ParseError struct {
	File    string
	Line    int
	Struct  string
	Field   string
	Message string
}

func (e ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s:%d: поле %s.%s: %s",
			e.File, e.Line, e.Struct, e.Field, e.Message)
	}
	if e.Struct != "" {
		return fmt.Sprintf("%s:%d: структура %s: %s",
			e.File, e.Line, e.Struct, e.Message)
	}
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}

// DirectiveError represents an error in directive validation
type DirectiveError struct {
	File      string
	Line      int
	Struct    string
	Field     string
	Directive string
	Message   string
	Severity  Severity
}

func (e DirectiveError) Error() string {
	return fmt.Sprintf("%s:%d: %s: поле %s.%s: директива %s: %s",
		e.File, e.Line, e.Severity, e.Struct, e.Field, e.Directive, e.Message)
}
