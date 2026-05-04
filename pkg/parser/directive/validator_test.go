package directive

import (
	"testing"

	"github.com/arkannsk/elval/pkg/errs"
	ann "github.com/arkannsk/elval/pkg/parser/annotations"
	"github.com/stretchr/testify/assert"
)

func TestValidate_UnknownDirective(t *testing.T) {
	info := FieldInfo{
		TypeName:   "string",
		StructName: "User",
		FieldName:  "Phone",
	}
	dir := ann.Directive{Type: "phone", Params: []string{}}
	loc := errs.Location{File: "user.go", Line: 24}

	diags := Validate(info, dir, loc)

	assert.Len(t, diags, 1)
	assert.Equal(t, errs.SeverityWarning, diags[0].Severity)
	assert.Equal(t, "validator", diags[0].Component)
	assert.Contains(t, diags[0].Message, "unknown directive 'phone'")
	assert.NotEmpty(t, diags[0].Suggestion, "should have suggestion with available directives")
	assert.False(t, HasErrors(diags), "unknown directive should not block generation")
}

func TestValidate_TypeMismatch_Error(t *testing.T) {
	info := FieldInfo{
		TypeName:   "int",
		StructName: "User",
		FieldName:  "Score",
	}
	dir := ann.Directive{Type: "pattern", Params: []string{"email"}}
	loc := errs.Location{File: "user.go", Line: 31}

	diags := Validate(info, dir, loc)

	assert.Len(t, diags, 1)
	assert.Equal(t, errs.SeverityError, diags[0].Severity)
	assert.Contains(t, diags[0].Message, "not supported for type int")
	assert.True(t, HasErrors(diags), "type mismatch should block directive")
}

func TestValidate_MissingParam_Warning(t *testing.T) {
	info := FieldInfo{
		TypeName:   "int",
		StructName: "User",
		FieldName:  "Age",
	}
	dir := ann.Directive{Type: "min", Params: []string{}}
	loc := errs.Location{File: "user.go", Line: 45}

	diags := Validate(info, dir, loc)

	assert.Len(t, diags, 1)
	assert.Equal(t, errs.SeverityWarning, diags[0].Severity)
	assert.Contains(t, diags[0].Message, "requires one parameter")
	assert.Contains(t, diags[0].Suggestion, "min:value", "should suggest usage example")
	assert.False(t, HasErrors(diags), "missing param is warning, not error")
}

func TestValidate_ValidDirective_Empty(t *testing.T) {
	info := FieldInfo{
		TypeName:   "string",
		StructName: "User",
		FieldName:  "Email",
	}
	dir := ann.Directive{Type: "required", Params: []string{}}
	loc := errs.Location{File: "user.go", Line: 10}

	diags := Validate(info, dir, loc)

	assert.Empty(t, diags, "valid directive should produce no diagnostics")
}

func TestValidate_Deprecated_Warning(t *testing.T) {
	info := FieldInfo{
		TypeName:   "int",
		StructName: "User",
		FieldName:  "Range",
	}
	dir := ann.Directive{Type: "min-max", Params: []string{"3", "50"}}
	loc := errs.Location{File: "user.go", Line: 50}

	diags := Validate(info, dir, loc)

	// Deprecated directive: warning + valid result (no error)
	assert.NotEmpty(t, diags)
	hasDeprecation := false
	for _, d := range diags {
		if d.Severity == errs.SeverityWarning &&
			d.Message == "directive 'min-max' is deprecated: Value range constraint — deprecated, use min and max separately" {
			hasDeprecation = true
			assert.Contains(t, d.Suggestion, "min-max:3,50")
		}
	}
	assert.True(t, hasDeprecation, "should warn about deprecation")
	assert.False(t, HasErrors(diags), "deprecation is warning, not error")
}

func TestValidate_CustomDirective_Skipped(t *testing.T) {
	info := FieldInfo{TypeName: "string", StructName: "User", FieldName: "Custom"}
	dir := ann.Directive{Type: "x-custom", Params: []string{"value"}}
	loc := errs.Location{File: "user.go", Line: 1}

	diags := Validate(info, dir, loc)

	assert.Empty(t, diags, "x-* directives should skip validation")
}
