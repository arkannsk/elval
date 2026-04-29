# Utils Package

Package utils provides helper functions for ElVal code generation.

## Functions

## Usage

#### func  BaseType

```go
func BaseType(ft parser.FieldType) string
```
BaseType returns the base name of the type, removing pointers and slices. If a
BaseType alias is specified, it returns that instead.

#### func  BuildGlobalRef

```go
func BuildGlobalRef(typeName, pkgPath, module string) string
```
BuildGlobalRef forms a global reference to a type, considering the structure of
the name.

#### func  Contains

```go
func Contains(s, substr string) bool
```
Contains reports whether substr is within s.

#### func  CountBodyFields

```go
func CountBodyFields(fields []parser.Field) int
```
CountBodyFields counts the number of fields that will be included in the OpenAPI
Schema (i.e., fields that are not HTTP parameters).

#### func  Dict

```go
func Dict(values ...any) (map[string]any, error)
```
Dict creates a map[string]any{} from a sequence of key-value pairs. It returns
an error if the number of arguments is odd or if keys are not strings.

#### func  FirstArg

```go
func FirstArg(args []parser.FieldType) parser.FieldType
```
FirstArg returns the first element of the FieldType slice, or an empty FieldType
if the slice is empty.

#### func  GenerateDecoratorCode

```go
func GenerateDecoratorCode(decType string, paramName string, fieldName string) string
```
GenerateDecoratorCode generates Go code for a specific decorator type. It takes
the decorator type (e.g., "ctx-get", "env-get"), the parameter name (if
applicable, e.g., key for ctx-get), and the target field name. It returns a
string containing the generated Go code snippet.

Supported decorator types:

    - ctx-get: Retrieves value from context.
    - httpctx-get: Retrieves header from HTTP request in context.
    - env-get: Retrieves value from environment variable.
    - time-now: Sets current time.
    - uuid-gen: Generates a new UUID.
    - trim, lower, upper: String manipulation functions.

#### func  GenerateParseCode

```go
func GenerateParseCode(fieldName string, goType string, sourceExpr string, isSlice bool) string
```
GenerateParseCode generates Go code for parsing an HTTP parameter into a
specific Go type. It handles primitives (int, float, bool, time.Time) and slices
of these types.

Parameters:

    - fieldName: The name of the struct field to assign the parsed value to.
    - goType: The base Go type (e.g., "int", "string", "time.Time").
    - sourceExpr: The expression to get the raw string value (e.g., 'r.URL.Query()["page"][0]').
    - isSlice: If true, generates code for parsing a slice; otherwise, parses a single value.

The function returns a string containing the complete Go code block, including
error handling using errs.NewParseRequestError. For slices, it iterates over the
source values. For primitives, it performs conversion and assigns the result to
v.fieldName.

#### func  GenerateSchemaTypeCode

```go
func GenerateSchemaTypeCode(typeName, prefix string) string
```
GenerateSchemaTypeCode generates assignment code for the OpenAPI schema type.

#### func  GetFieldAnnotationValue

```go
func GetFieldAnnotationValue(field parser.Field, annotationType string) string
```
GetFieldAnnotationValue extracts the value of an annotation of the specified
type from the field.

#### func  GlobalRefFor

```go
func GlobalRefFor(typeName, typePkgPath, typeMod, structPkgPath, structMod string) string
```
GlobalRefFor forms a global reference to a type.

#### func  HasDirective

```go
func HasDirective(directives []ann.Directive, name string) bool
```
HasDirective checks if a directive with the specified name exists in the list.

#### func  HasOptional

```go
func HasOptional(directives []ann.Directive) bool
```
HasOptional checks if the 'optional' directive is present in the list of
directives.

#### func  HasPrefix

```go
func HasPrefix(s, prefix string) bool
```
HasPrefix tests whether the string s begins with prefix.

#### func  HasSuffix

```go
func HasSuffix(s, suffix string) bool
```
HasSuffix tests whether the string s ends with suffix.

#### func  HasTime

```go
func HasTime(structs []parser.Struct) bool
```
HasTime checks if the list of structs contains fields of type time.Time or
time.Duration.

#### func  IsCustomDirective

```go
func IsCustomDirective(dirType string) bool
```
IsCustomDirective checks if the directive type is custom (starts with x-).

#### func  IsFieldRequired

```go
func IsFieldRequired(field parser.Field) bool
```
IsFieldRequired checks if the field is required (presence of 'required'
directive).

#### func  IsNumericOrBoolOrTime

```go
func IsNumericOrBoolOrTime(goType string) bool
```
IsNumericOrBoolOrTime checks if the type is numeric, boolean, or time-related.

#### func  IsPrimitive

```go
func IsPrimitive(name string, primitives map[string]bool) bool
```
IsPrimitive checks if the type name is a primitive Go type. It requires a map of
known primitives to be passed as an argument.

#### func  Itoa

```go
func Itoa(i int) string
```
Itoa converts an int to a string.

#### func  OaString

```go
func OaString(val string) string
```
OaString wraps the value in a Go string literal (with quotes).

#### func  RegexMatch

```go
func RegexMatch(re, s string) bool
```
RegexMatch reports whether the string s matches the regular expression re.

#### func  SafeExample

```go
func SafeExample(val string) string
```
SafeExample formats a value for use as an example in OpenAPI. JSON objects and
arrays are wrapped in backticks.

#### func  Split

```go
func Split(s, sep string) []string
```
Split slices s into all substrings separated by sep and returns a slice of the
substrings between those separators.

#### func  SplitList

```go
func SplitList(s, sep string) []string
```
SplitList splits a string by a separator and trims whitespace from each part.

#### func  Title

```go
func Title(s string) string
```
Title returns a copy of the string s with all Unicode letters that begin words
mapped to their title case.

#### func  ToLower

```go
func ToLower(s string) string
```
ToLower returns the string s converted to lowercase.

#### func  ToOpenAPIType

```go
func ToOpenAPIType(goType string) string
```
ToOpenAPIType converts a Go type name to an OpenAPI type.

#### func  Trim

```go
func Trim(s string) string
```
Trim removes leading and trailing whitespace from s.

#### func  TrimBrackets

```go
func TrimBrackets(s string) string
```
TrimBrackets removes square brackets '[' and ']' from the string.

#### func  TrimPrefix

```go
func TrimPrefix(s, prefix string) string
```
TrimPrefix returns s without the provided leading prefix string.

#### func  TrimQuotes

```go
func TrimQuotes(s string) string
```
TrimQuotes removes outer quotes (" or ') from the string s.

#### func  TrimStar

```go
func TrimStar(s string) string
```
TrimStar removes the '*' character from the beginning of the string.

#### func  TrimSuffix

```go
func TrimSuffix(s, suffix string) string
```
TrimSuffix returns s without the provided trailing suffix string.

#### func  UniqueRef

```go
func UniqueRef(typeName, typePackage, typeModule, structPackage, structModule string) string
```
UniqueRef forms a unique reference to a type.
