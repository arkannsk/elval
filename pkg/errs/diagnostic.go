package errs

import "fmt"

// Severity уровень важности
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Location позиция в исходном коде
type Location struct {
	File   string
	Line   int
	Column int
}

// Diagnostic сообщение от валидатора
type Diagnostic struct {
	Loc        Location
	Severity   Severity
	Component  string
	Directive  string
	StructName string
	FieldName  string
	Message    string
	Suggestion string
}

func (d Diagnostic) String() string {
	loc := d.Loc.File
	if d.Loc.Line > 0 {
		loc = fmt.Sprintf("%s:%d", loc, d.Loc.Line)
	}

	ctx := ""
	if d.StructName != "" {
		ctx = d.StructName
		if d.FieldName != "" {
			ctx += "." + d.FieldName
		}
		ctx += ": "
	}

	return fmt.Sprintf("[%s] %s (%s)%s %s",
		d.Severity, loc, d.Component, ctx, d.Message)
}

func NewWarning(loc Location, directive, structName, fieldName, msg string) Diagnostic {
	return Diagnostic{
		Loc: loc, Severity: SeverityWarning, Component: "validator",
		Directive: directive, StructName: structName, FieldName: fieldName,
		Message: msg,
	}
}

func NewError(loc Location, directive, structName, fieldName, msg string) Diagnostic {
	return Diagnostic{
		Loc: loc, Severity: SeverityError, Component: "validator",
		Directive: directive, StructName: structName, FieldName: fieldName,
		Message: msg,
	}
}

func (d Diagnostic) WithSuggestion(s string) Diagnostic {
	d.Suggestion = s
	return d
}
