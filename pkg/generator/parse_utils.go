package generator

import "fmt"

// GenerateParseCode генерирует код Go для парсинга HTTP параметров в нужный тип, пока удобнее чем в шаблонах
func GenerateParseCode(fieldName string, goType string, sourceExpr string, isSlice bool) string {
	if isSlice {
		return generateSliceParseCode(fieldName, goType, sourceExpr)
	}
	return generatePrimitiveParseCode(fieldName, goType, sourceExpr)
}

func generatePrimitiveParseCode(fieldName, goType, valExpr string) string {
	switch goType {
	case "string":
		return fmt.Sprintf("v.%s = %s", fieldName, valExpr)

	case "int", "int8", "int16", "int32", "int64":
		parsedVar := "parsedInt"
		code := fmt.Sprintf(`%s, err := strconv.ParseInt(%s, 10, 64)
if err != nil {
	return errs.NewParseRequestError(%q, %s, "invalid integer")
}`, parsedVar, valExpr, fieldName, valExpr)

		if goType == "int64" {
			return code + fmt.Sprintf("\nv.%s = %s", fieldName, parsedVar)
		}
		return code + fmt.Sprintf("\nv.%s = %s(%s)", fieldName, goType, parsedVar)

	case "uint", "uint8", "uint16", "uint32", "uint64":
		parsedVar := "parsedUint"
		code := fmt.Sprintf(`%s, err := strconv.ParseUint(%s, 10, 64)
if err != nil {
	return errs.NewParseRequestError(%q, %s, "invalid uint")
}`, parsedVar, valExpr, fieldName, valExpr)

		if goType == "uint64" {
			return code + fmt.Sprintf("\nv.%s = %s", fieldName, parsedVar)
		}
		return code + fmt.Sprintf("\nv.%s = %s(%s)", fieldName, goType, parsedVar)

	case "float32", "float64":
		parsedVar := "parsedFloat"
		code := fmt.Sprintf(`%s, err := strconv.ParseFloat(%s, 64)
if err != nil {
	return errs.NewParseRequestError(%q, %s, "invalid float")
}`, parsedVar, valExpr, fieldName, valExpr)

		if goType == "float64" {
			return code + fmt.Sprintf("\nv.%s = %s", fieldName, parsedVar)
		}
		return code + fmt.Sprintf("\nv.%s = %s(%s)", fieldName, goType, parsedVar)

	case "bool":
		parsedVar := "parsedBool"
		code := fmt.Sprintf(`%s, err := strconv.ParseBool(%s)
if err != nil {
	return errs.NewParseRequestError(%q, %s, "invalid bool")
}`, parsedVar, valExpr, fieldName, valExpr)
		return code + fmt.Sprintf("\nv.%s = %s", fieldName, parsedVar)

	case "time.Time":
		parsedVar := "parsedTime"
		code := fmt.Sprintf(`%s, err := time.Parse(time.RFC3339, %s)
if err != nil {
	return errs.NewParseRequestError(%q, %s, "invalid time format (RFC3339)")
}`, parsedVar, valExpr, fieldName, valExpr)
		return code + fmt.Sprintf("\nv.%s = %s", fieldName, parsedVar)

	default:
		return fmt.Sprintf("// Unparsed type %s from %s\nv.%s = %s", goType, valExpr, fieldName, valExpr)
	}
}

func generateSliceParseCode(fieldName, baseType, valsExpr string) string {
	switch baseType {
	case "string":
		return fmt.Sprintf("v.%s = %s", fieldName, valsExpr)

	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf(`var result []%s
for _, v := range %s {
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return errs.NewParseRequestError(%q, v, "invalid integer in slice")
	}
	result = append(result, %s(val))
}
v.%s = result`, baseType, valsExpr, fieldName, baseType, fieldName)

	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf(`var result []%s
for _, v := range %s {
	val, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return errs.NewParseRequestError(%q, v, "invalid uint in slice")
	}
	result = append(result, %s(val))
}
v.%s = result`, baseType, valsExpr, fieldName, baseType, fieldName)

	case "float32", "float64":
		return fmt.Sprintf(`var result []%s
for _, v := range %s {
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return errs.NewParseRequestError(%q, v, "invalid float in slice")
	}
	result = append(result, %s(val))
}
v.%s = result`, baseType, valsExpr, fieldName, baseType, fieldName)

	case "bool":
		return fmt.Sprintf(`var result []%s
for _, v := range %s {
	val, err := strconv.ParseBool(v)
	if err != nil {
		return errs.NewParseRequestError(%q, v, "invalid bool in slice")
	}
	result = append(result, val)
}
v.%s = result`, baseType, valsExpr, fieldName, fieldName)

	default:
		return fmt.Sprintf("// Unparsed slice of %s from %s\nv.%s = %s", baseType, valsExpr, fieldName, valsExpr)
	}
}
