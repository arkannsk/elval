// Package utils provides helper functions for ElVal code generation.
package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/arkannsk/elval/pkg/parser"
	ann "github.com/arkannsk/elval/pkg/parser/annotations"
)

// HasTime checks if the list of structs contains fields of type time.Time or time.Duration.
func HasTime(structs []parser.Struct) bool {
	for _, s := range structs {
		for _, f := range s.Fields {
			if f.Type.Name == "time.Time" || f.Type.Name == "time.Duration" {
				return true
			}
		}
	}
	return false
}

// Dict creates a map[string]any{} from a sequence of key-value pairs.
// It returns an error if the number of arguments is odd or if keys are not strings.
func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// HasOptional checks if the 'optional' directive is present in the list of directives.
func HasOptional(directives []ann.Directive) bool {
	for _, d := range directives {
		if d.Type == "optional" {
			return true
		}
	}
	return false
}

// Itoa converts an int to a string.
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// IsCustomDirective checks if the directive type is custom (starts with x-).
func IsCustomDirective(dirType string) bool {
	return strings.HasPrefix(dirType, "x-")
}

// ToLower returns the string s converted to lowercase.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Contains reports whether substr is within s.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HasSuffix tests whether the string s ends with suffix.
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// HasPrefix tests whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// RegexMatch reports whether the string s matches the regular expression re.
func RegexMatch(re, s string) bool {
	matched, _ := regexp.MatchString(re, s)
	return matched
}

// Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators.
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Title returns a copy of the string s with all Unicode letters that begin words mapped to their title case.
func Title(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Trim removes leading and trailing whitespace from s.
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimPrefix returns s without the provided leading prefix string.
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimSuffix returns s without the provided trailing suffix string.
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// TrimQuotes removes outer quotes (" or ') from the string s.
func TrimQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// OaString wraps the value in a Go string literal (with quotes).
func OaString(val string) string {
	val = TrimQuotes(val)
	return fmt.Sprintf("%q", val)
}

// TrimStar removes the '*' character from the beginning of the string.
func TrimStar(s string) string {
	return strings.TrimPrefix(s, "*")
}

// FirstArg returns the first element of the FieldType slice, or an empty FieldType if the slice is empty.
func FirstArg(args []parser.FieldType) parser.FieldType {
	if len(args) > 0 {
		return args[0]
	}
	return parser.FieldType{}
}

// IsPrimitive checks if the type name is a primitive Go type.
// It requires a map of known primitives to be passed as an argument.
func IsPrimitive(name string, primitives map[string]bool) bool {
	_, ok := primitives[name]
	return ok
}

// HasDirective checks if a directive with the specified name exists in the list.
func HasDirective(directives []ann.Directive, name string) bool {
	for _, d := range directives {
		if d.Type == name {
			return true
		}
	}
	return false
}

// TrimBrackets removes square brackets '[' and ']' from the string.
func TrimBrackets(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

// BaseType returns the base name of the type, removing pointers and slices.
// If a BaseType alias is specified, it returns that instead.
func BaseType(ft parser.FieldType) string {
	name := ft.Name
	name = strings.TrimPrefix(name, "*")
	if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
		name = strings.TrimPrefix(strings.TrimSuffix(name, "]"), "[")
	}
	if ft.BaseType != "" {
		return ft.BaseType
	}
	return name
}

// UniqueRef forms a unique reference to a type.
func UniqueRef(typeName, typePackage, typeModule, structPackage, structModule string) string {
	if typeModule != "" && typeModule != structModule {
		return fmt.Sprintf("%s/%s.%s", typeModule, typePackage, typeName)
	}
	if typePackage != "" && typePackage != structPackage {
		return fmt.Sprintf("%s.%s", typePackage, typeName)
	}
	return typeName
}

// SafeExample formats a value for use as an example in OpenAPI.
// JSON objects and arrays are wrapped in backticks.
func SafeExample(val string) string {
	val = strings.TrimSpace(val)
	if len(val) > 0 && (val[0] == '{' || val[0] == '[') {
		return fmt.Sprintf("`%s`", val)
	}
	return val
}

// GlobalRefFor forms a global reference to a type.
func GlobalRefFor(typeName, typePkgPath, typeMod, structPkgPath, structMod string) string {
	if typeMod != "" && typeMod != structMod {
		return fmt.Sprintf("%s/%s.%s", typeMod, typePkgPath, typeName)
	}
	if typePkgPath != "" && typePkgPath != structPkgPath {
		return fmt.Sprintf("%s.%s", typePkgPath, typeName)
	}
	return typeName
}

// BuildGlobalRef forms a global reference to a type, considering the structure of the name.
func BuildGlobalRef(typeName, pkgPath, module string) string {
	if strings.Contains(typeName, "/") {
		return fmt.Sprintf("%s/%s", module, typeName)
	}
	if strings.Contains(typeName, ".") {
		return fmt.Sprintf("%s/%s/%s", module, pkgPath, typeName)
	}
	return fmt.Sprintf("%s/%s.%s", module, pkgPath, typeName)
}

// SplitList splits a string by a separator and trims whitespace from each part.
func SplitList(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		result = append(result, strings.TrimSpace(p))
	}
	return result
}

// GenerateSchemaTypeCode generates assignment code for the OpenAPI schema type.
func GenerateSchemaTypeCode(typeName, prefix string) string {
	switch typeName {
	case "string":
		return fmt.Sprintf("%s.Type = \"string\"", prefix)
	case "time.Time":
		return fmt.Sprintf(`%s.Type = "string" %s.Format = "date-time"`, prefix, prefix)
	case "bool":
		return fmt.Sprintf("%s.Type = \"boolean\"", prefix)
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s.Type = \"integer\"", prefix)
	case "float32", "float64":
		return fmt.Sprintf("%s.Type = \"number\"", prefix)
	default:
		return ""
	}
}

// IsNumericOrBoolOrTime checks if the type is numeric, boolean, or time-related.
func IsNumericOrBoolOrTime(goType string) bool {
	switch goType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"bool",
		"time.Time", "time.Duration":
		return true
	default:
		return false
	}
}

// ToOpenAPIType converts a Go type name to an OpenAPI type.
func ToOpenAPIType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string" // In OpenAPI, time is usually string with format date-time
	default:
		if strings.HasPrefix(goType, "[") || strings.HasPrefix(goType, "*") {
			return "array"
		}
		return "object"
	}
}

// GetFieldAnnotationValue extracts the value of an annotation of the specified type from the field.
func GetFieldAnnotationValue(field parser.Field, annotationType string) string {
	for _, ann := range field.OaAnnotations {
		if ann.Type == annotationType {
			return ann.Value
		}
	}
	return ""
}

// IsFieldRequired checks if the field is required (presence of 'required' directive).
func IsFieldRequired(field parser.Field) bool {
	for _, dir := range field.Directives {
		if dir.Type == "required" {
			return true
		}
	}
	return false
}

// CountBodyFields counts the number of fields that will be included in the OpenAPI Schema
// (i.e., fields that are not HTTP parameters).
func CountBodyFields(fields []parser.Field) int {
	count := 0
	for _, f := range fields {
		if f.OaIn == "" {
			count++
		}
	}
	return count
}

func JoinQuoted(params []string) string {
	quoted := make([]string, len(params))
	for i, p := range params {
		quoted[i] = fmt.Sprintf("%q", p)
	}
	return strings.Join(quoted, ", ")
}

func JoinUnquoted(params []string) string {
	return strings.Join(params, ", ")
}
