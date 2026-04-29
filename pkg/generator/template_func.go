package generator

import (
	"text/template"

	"github.com/arkannsk/elval/pkg/generator/utils"
)

var templateFucMap = template.FuncMap{
	"hasTime":           utils.HasTime,
	"dict":              utils.Dict,
	"hasOptional":       utils.HasOptional,
	"itoa":              utils.Itoa,
	"isCustomDirective": utils.IsCustomDirective,
	"toLower":           utils.ToLower,
	"contains":          utils.Contains,
	"hasSuffix":         utils.HasSuffix,
	"hasPrefix":         utils.HasPrefix,
	"regexMatch":        utils.RegexMatch,
	"split":             utils.Split,
	"title":             utils.Title,
	"trim":              utils.Trim,
	"trimPrefix":        utils.TrimPrefix,
	"trimSuffix":        utils.TrimSuffix,
	"trimQuotes":        utils.TrimQuotes,
	"oaString":          utils.OaString,
	"trimStar":          utils.TrimStar,
	"firstArg":          utils.FirstArg,
	"isPrimitive": func(name string) bool {
		_, ok := primitives[name]
		return ok
	},
	"hasDirective":            utils.HasDirective,
	"trimBrackets":            utils.TrimBrackets,
	"baseType":                utils.BaseType,
	"uniqueRef":               utils.UniqueRef,
	"safeExample":             utils.SafeExample,
	"globalRefFor":            utils.GlobalRefFor,
	"buildGlobalRef":          utils.BuildGlobalRef,
	"splitList":               utils.SplitList,
	"generateSchemaTypeCode":  utils.GenerateSchemaTypeCode,
	"GenerateDecoratorCode":   utils.GenerateDecoratorCode,
	"GenerateParseCode":       utils.GenerateParseCode,
	"IsNumericOrBoolOrTime":   utils.IsNumericOrBoolOrTime,
	"ToOpenAPIType":           utils.ToOpenAPIType,
	"GetFieldAnnotationValue": utils.GetFieldAnnotationValue,
	"IsFieldRequired":         utils.IsFieldRequired,
	"CountBodyFields":         utils.CountBodyFields,
}
