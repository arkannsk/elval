package generator

import (
	"fmt"
	"strings"
)

// GenerateParseCode генерирует код Go для парсинга HTTP параметров в нужный тип
func GenerateParseCode(fieldName string, goType string, sourceExpr string, isSlice bool) string {
	if isSlice {
		return generateSliceParseCode(fieldName, goType, sourceExpr)
	}
	return generatePrimitiveParseCode(fieldName, goType, sourceExpr)
}

func generatePrimitiveParseCode(fieldName, goType, valExpr string) string {
	var sb strings.Builder

	switch goType {
	case "string":
		sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, valExpr))

	case "int", "int8", "int16", "int32", "int64":
		parsedVar := "parsedInt"
		sb.WriteString(fmt.Sprintf("\t%s, err := strconv.ParseInt(%s, 10, 64)\n", parsedVar, valExpr))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\treturn errs.NewParseRequestError(%q, %s, \"invalid integer\")\n", fieldName, valExpr))
		sb.WriteString("\t}\n")
		if goType == "int64" {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, parsedVar))
		} else {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s(%s)\n", fieldName, goType, parsedVar))
		}

	case "uint", "uint8", "uint16", "uint32", "uint64":
		parsedVar := "parsedUint"
		sb.WriteString(fmt.Sprintf("\t%s, err := strconv.ParseUint(%s, 10, 64)\n", parsedVar, valExpr))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\treturn errs.NewParseRequestError(%q, %s, \"invalid uint\")\n", fieldName, valExpr))
		sb.WriteString("\t}\n")
		if goType == "uint64" {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, parsedVar))
		} else {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s(%s)\n", fieldName, goType, parsedVar))
		}

	case "float32", "float64":
		parsedVar := "parsedFloat"
		sb.WriteString(fmt.Sprintf("\t%s, err := strconv.ParseFloat(%s, 64)\n", parsedVar, valExpr))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\treturn errs.NewParseRequestError(%q, %s, \"invalid float\")\n", fieldName, valExpr))
		sb.WriteString("\t}\n")
		if goType == "float64" {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, parsedVar))
		} else {
			sb.WriteString(fmt.Sprintf("\tv.%s = %s(%s)\n", fieldName, goType, parsedVar))
		}

	case "bool":
		parsedVar := "parsedBool"
		sb.WriteString(fmt.Sprintf("\t%s, err := strconv.ParseBool(%s)\n", parsedVar, valExpr))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\treturn errs.NewParseRequestError(%q, %s, \"invalid bool\")\n", fieldName, valExpr))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, parsedVar))

	case "time.Time":
		parsedVar := "parsedTime"
		sb.WriteString(fmt.Sprintf("\t%s, err := time.Parse(time.RFC3339, %s)\n", parsedVar, valExpr))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\treturn errs.NewParseRequestError(%q, %s, \"invalid time format (RFC3339)\")\n", fieldName, valExpr))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, parsedVar))

	default:
		sb.WriteString(fmt.Sprintf("\t// TODO: Parse type %s from %s\n", goType, valExpr))
		sb.WriteString(fmt.Sprintf("\tv.%s = %s // Unparsed\n", fieldName, valExpr))
	}

	return sb.String()
}

func generateSliceParseCode(fieldName, baseType, valsExpr string) string {
	var sb strings.Builder

	switch baseType {
	case "string":
		sb.WriteString(fmt.Sprintf("\tv.%s = %s\n", fieldName, valsExpr))

	case "int", "int8", "int16", "int32", "int64":
		sb.WriteString(fmt.Sprintf("\tvar result []%s\n", baseType))
		sb.WriteString(fmt.Sprintf("\tfor _, v := range %s {\n", valsExpr))
		sb.WriteString("\t\tval, err := strconv.ParseInt(v, 10, 64)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn errs.NewParseRequestError(%q, v, \"invalid integer in slice\")\n", fieldName))
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\tresult = append(result, %s(val))\n", baseType))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = result\n", fieldName))

	case "uint", "uint8", "uint16", "uint32", "uint64":
		sb.WriteString(fmt.Sprintf("\tvar result []%s\n", baseType))
		sb.WriteString(fmt.Sprintf("\tfor _, v := range %s {\n", valsExpr))
		sb.WriteString("\t\tval, err := strconv.ParseUint(v, 10, 64)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn errs.NewParseRequestError(%q, v, \"invalid uint in slice\")\n", fieldName))
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\tresult = append(result, %s(val))\n", baseType))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = result\n", fieldName))

	case "float32", "float64":
		sb.WriteString(fmt.Sprintf("\tvar result []%s\n", baseType))
		sb.WriteString(fmt.Sprintf("\tfor _, v := range %s {\n", valsExpr))
		sb.WriteString("\t\tval, err := strconv.ParseFloat(v, 64)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn errs.NewParseRequestError(%q, v, \"invalid float in slice\")\n", fieldName))
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\tresult = append(result, %s(val))\n", baseType))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = result\n", fieldName))

	case "bool":
		sb.WriteString(fmt.Sprintf("\tvar result []%s\n", baseType))
		sb.WriteString(fmt.Sprintf("\tfor _, v := range %s {\n", valsExpr))
		sb.WriteString("\t\tval, err := strconv.ParseBool(v)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn errs.NewParseRequestError(%q, v, \"invalid bool in slice\")\n", fieldName))
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tresult = append(result, val)\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tv.%s = result\n", fieldName))

	default:
		sb.WriteString(fmt.Sprintf("\t// TODO: Parse slice of %s from %s\n", baseType, valsExpr))
		sb.WriteString(fmt.Sprintf("\tv.%s = %s // Unparsed\n", fieldName, valsExpr))
	}

	return sb.String()
}
