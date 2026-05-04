// pkg/parser/directive/validators.go
package directive

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
)

// ValidatorFunc validates a single directive.
// dirType is passed explicitly so validators can use it in messages.
type ValidatorFunc func(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic

// Registry maps directive types to their validator functions.
var Registry = map[Type]ValidatorFunc{
	TypeRequired:   validateNoParams,
	TypeOptional:   validateNoParams,
	TypeMin:        validateNumericLimit,
	TypeMax:        validateNumericLimit,
	TypeLen:        validateLength,
	TypeMinMax:     validateMinMaxRange,
	TypePattern:    validatePattern,
	TypeEnum:       validateEnum,
	TypeURL:        validateNoParams,
	TypeHTTPURL:    validateNoParams,
	TypeDSN:        validateNoParams,
	TypeContains:   validateStringParam,
	TypeStartsWith: validateStringParam,
	TypeEndsWith:   validateStringParam,
	TypeNotZero:    validateNoParams,
	TypeBefore:     validateDateCompare,
	TypeAfter:      validateDateCompare,
	TypeDate:       validateDateFormats,
	TypeEq:         validateComparisonParam,
	TypeNeq:        validateComparisonParam,
	TypeLt:         validateComparisonParam,
	TypeLte:        validateComparisonParam,
	TypeGt:         validateComparisonParam,
	TypeGte:        validateComparisonParam,
	TypeRequiredIf: validateRequiredIf,
}

// ─────────────────────────────────────────────────────────────
// Common helpers
// ─────────────────────────────────────────────────────────────

func warn(loc errs.Location, dir, structName, fieldName, msg, suggestion string) errs.Diagnostic {
	return errs.NewWarning(loc, dir, structName, fieldName, msg).WithSuggestion(suggestion)
}

func err(loc errs.Location, dir, structName, fieldName, msg string) errs.Diagnostic {
	return errs.NewError(loc, dir, structName, fieldName, msg)
}

// checkParamCount validates the number of parameters for a directive.
func checkParamCount(expected, actual int, dirType Type, info FieldInfo, loc errs.Location) []errs.Diagnostic {
	dName := string(dirType)

	if expected == ParamNone && actual > 0 {
		return []errs.Diagnostic{warn(loc, dName, info.StructName, info.FieldName,
			fmt.Sprintf("directive '%s' does not accept parameters", dName),
			fmt.Sprintf("usage: @evl:validate %s", dName))}
	}
	if expected == ParamOne && actual == 0 {
		return []errs.Diagnostic{warn(loc, dName, info.StructName, info.FieldName,
			fmt.Sprintf("directive '%s' requires one parameter", dName),
			fmt.Sprintf("example: @evl:validate %s:value", dName))}
	}
	if expected == ParamTwo && actual != 2 {
		return []errs.Diagnostic{warn(loc, dName, info.StructName, info.FieldName,
			fmt.Sprintf("directive '%s' requires two parameters", dName),
			fmt.Sprintf("example: @evl:validate %s:min,max", dName))}
	}
	return nil
}

// checkTypeAllowed validates that the directive is applicable to the field type.
func checkTypeAllowed(info FieldInfo, allowed []string, dirType Type, loc errs.Location) []errs.Diagnostic {
	baseType := info.TypeName
	if info.IsPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}

	if info.IsGeneric && len(info.GenericArgs) > 0 {
		for _, arg := range info.GenericArgs {
			if diags := checkTypeAllowedForSingleType(arg, allowed, dirType, loc); len(diags) > 0 {
				return diags
			}
		}
		return nil
	}

	return checkTypeAllowedForSingleType(info, allowed, dirType, loc)
}

func checkTypeAllowedForSingleType(info FieldInfo, allowed []string, dirType Type, loc errs.Location) []errs.Diagnostic {
	typeName := info.TypeName
	if info.IsPointer {
		typeName = strings.TrimPrefix(typeName, "*")
	}

	baseType := typeName
	if info.BaseType != "" && info.BaseType != typeName {
		baseType = info.BaseType
	}

	actualType := baseType
	if info.IsSlice {
		actualType = "slice"
	} else if info.IsPointer {
		actualType = "pointer"
	}

	for _, t := range allowed {
		if t == "any" {
			return nil
		}
		if t == typeName || t == baseType || t == actualType {
			return nil
		}
		if info.IsPointer && info.IsStruct && t == "pointer" {
			return nil
		}
		if strings.HasSuffix(typeName, "."+t) || strings.HasSuffix(baseType, "."+t) {
			return nil
		}
	}

	msg := fmt.Sprintf("directive '%s' is not supported for type %s", dirType, info.TypeName)
	if info.BaseType != "" && info.BaseType != info.TypeName {
		msg += fmt.Sprintf(" (alias for %s)", info.BaseType)
	}

	return []errs.Diagnostic{err(loc, string(dirType), info.StructName, info.FieldName, msg)}
}

func parseIntParam(param string) (int, error) {
	val, err := strconv.Atoi(param)
	if err != nil || val < 0 {
		return 0, fmt.Errorf("must be non-negative integer")
	}
	return val, nil
}

func parseFloatParam(param string) (float64, error) {
	return strconv.ParseFloat(param, 64)
}

func parseDurationParam(param string) error {
	_, err := time.ParseDuration(param)
	return err
}

func parseDateParam(param string) bool {
	formats := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if _, err := time.Parse(f, param); err == nil {
			return true
		}
	}
	return false
}

func validateNoParams(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)
	return checkParamCount(meta.ParamCount, len(params), dirType, info, loc)
}

func validateNumericLimit(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if diags := checkTypeAllowed(info, meta.AllowedTypes, dirType, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil // already warned by checkParamCount
	}

	param := params[0]

	// Определяем тип для парсинга параметра — учитываем алиасы и дженерики
	typeToCheck := info
	if info.IsGeneric && len(info.GenericArgs) > 0 {
		typeToCheck = info.GenericArgs[0]
	}

	// Для алиасов берём базовый тип
	baseType := typeToCheck.TypeName
	if typeToCheck.IsPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}
	if typeToCheck.BaseType != "" && typeToCheck.BaseType != baseType {
		baseType = typeToCheck.BaseType
	}

	// Теперь baseType = "string" для алиаса "type Theme string"
	// Handle time.Duration specially
	if strings.Contains(baseType, "Duration") {
		if err := parseDurationParam(param); err != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				fmt.Sprintf("parameter '%s' is not a valid duration", param),
				"examples: 1h, 30m, 500ms, 10s")}
		}
		return nil
	}

	// Handle strings and slices (length constraints)
	if typeToCheck.IsSlice || baseType == "string" {
		if _, err := parseIntParam(param); err != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				fmt.Sprintf("parameter '%s' must be a non-negative integer", param),
				"")}
		}
		return nil
	}

	// Handle numeric types
	if _, err := parseFloatParam(param); err != nil {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			fmt.Sprintf("parameter '%s' must be a number", param),
			"")}
	}
	return nil
}

func validateComparisonParam(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil
	}

	// For non-string, non-bool types, parameter must be numeric
	baseType := info.TypeName
	if info.IsPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}
	if baseType != "string" && baseType != "bool" {
		if _, err := parseFloatParam(params[0]); err != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				fmt.Sprintf("parameter '%s' must be a number for type %s", params[0], baseType),
				"")}
		}
	}
	return nil
}

func validatePattern(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if diags := checkTypeAllowed(info, meta.AllowedTypes, dirType, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil
	}
	if params[0] == "" {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"pattern parameter cannot be empty", "")}
	}

	// Check if using a predefined pattern name
	if presets := meta.PredefinedPatterns; presets != nil {
		if _, isPreset := presets[params[0]]; !isPreset {
			// Not a preset — assume it's a custom regex, but warn if it looks like a name
			if matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, params[0]); matched {
				keys := make([]string, 0, len(presets))
				for k := range presets {
					keys = append(keys, k)
				}
				return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
					fmt.Sprintf("unknown pattern preset '%s'", params[0]),
					fmt.Sprintf("available presets: %s", strings.Join(keys, ", ")))}
			}
		}
	}
	return nil
}

func validateStringParam(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	if diags := checkParamCount(ParamOne, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil
	}
	if params[0] == "" {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"parameter cannot be empty", "")}
	}
	return nil
}

func validateEnum(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"enum directive requires at least one allowed value",
			"example: @evl:validate enum:active,inactive,pending")}
	}
	return checkTypeAllowed(info, meta.AllowedTypes, dirType, loc)
}

func validateLength(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil
	}
	if _, err := parseIntParam(params[0]); err != nil {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			fmt.Sprintf("parameter '%s' must be a non-negative integer", params[0]),
			"")}
	}
	return checkTypeAllowed(info, meta.AllowedTypes, dirType, loc)
}

func validateMinMaxRange(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) != 2 {
		return nil // already warned
	}

	baseType := info.TypeName
	if info.IsPointer {
		baseType = strings.TrimPrefix(baseType, "*")
	}
	minVal, maxVal := params[0], params[1]

	// Handle time.Duration
	if strings.Contains(baseType, "Duration") {
		if err := parseDurationParam(minVal); err != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				"minimum value must be a valid duration", "example: 1h")}
		}
		if err := parseDurationParam(maxVal); err != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				"maximum value must be a valid duration", "example: 24h")}
		}
		return nil
	}

	// Handle strings/slices (length range)
	if info.IsSlice || baseType == "string" {
		minI, errMin := parseIntParam(minVal)
		maxI, errMax := parseIntParam(maxVal)
		if errMin != nil || errMax != nil {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				"parameters must be non-negative integers", "")}
		}
		if minI > maxI {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				fmt.Sprintf("min (%d) cannot be greater than max (%d)", minI, maxI), "")}
		}
		return nil
	}

	// Handle numeric types
	minF, errMin := parseFloatParam(minVal)
	maxF, errMax := parseFloatParam(maxVal)
	if errMin != nil || errMax != nil {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"parameters must be numbers", "")}
	}
	if minF > maxF {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			fmt.Sprintf("min (%v) cannot be greater than max (%v)", minF, maxF), "")}
	}
	return checkTypeAllowed(info, meta.AllowedTypes, dirType, loc)
}

func validateDateCompare(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if diags := checkTypeAllowed(info, meta.AllowedTypes, dirType, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return nil
	}
	if !parseDateParam(params[0]) {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			fmt.Sprintf("parameter '%s' is not a valid date", params[0]),
			"supported formats: 2006-01-02, RFC3339")}
	}
	return nil
}

func validateDateFormats(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) == 0 {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"date directive requires at least one format",
			"example: @evl:validate date:RFC3339,2006-01-02")}
	}

	// Validate each format parameter
	knownFormats := map[string]bool{
		"RFC3339": true, "RFC3339Nano": true,
		"2006-01-02": true, "15:04:05": true, "2006-01-02T15:04:05Z07:00": true,
		"2006-01-02 15:04:05": true,
	}
	for _, p := range params {
		if !knownFormats[p] && !strings.Contains(p, "-") && !strings.Contains(p, ":") {
			return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
				fmt.Sprintf("unknown date format '%s'", p),
				"examples: RFC3339, 2006-01-02, 15:04:05")}
		}
	}
	return nil
}

func validateRequiredIf(dirType Type, info FieldInfo, params []string, loc errs.Location) []errs.Diagnostic {
	meta := GetInfo(dirType)

	if diags := checkParamCount(meta.ParamCount, len(params), dirType, info, loc); len(diags) > 0 {
		return diags
	}
	if len(params) != 2 {
		return nil // already warned
	}
	if params[0] == "" {
		return []errs.Diagnostic{warn(loc, string(dirType), info.StructName, info.FieldName,
			"first parameter (field name) cannot be empty", "format: @evl:validate required_if:FieldName value")}
	}
	return nil
}
