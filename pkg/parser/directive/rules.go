package directive

type FieldInfo struct {
	TypeName    string
	IsSlice     bool
	IsPointer   bool
	IsGeneric   bool
	GenericArgs []FieldInfo
	BaseType    string
	IsStruct    bool
	StructName  string
	FieldName   string
}

type Info struct {
	Description        string
	AllowedTypes       []string
	ParamCount         int
	Example            string
	Deprecated         bool
	PredefinedPatterns map[string]string
}

const (
	ParamNone      = 0
	ParamOne       = 1
	ParamTwo       = 2
	ParamUnbounded = -1
)

// Type represents a validation directive type.
type Type string

// Directive type constants
const (
	TypeRequired   Type = "required"
	TypeOptional   Type = "optional"
	TypeMin        Type = "min"
	TypeMax        Type = "max"
	TypeLen        Type = "len"
	TypeMinMax     Type = "min-max"
	TypePattern    Type = "pattern"
	TypeEnum       Type = "enum"
	TypeURL        Type = "url"
	TypeHTTPURL    Type = "http_url"
	TypeDSN        Type = "dsn"
	TypeContains   Type = "contains"
	TypeStartsWith Type = "starts_with"
	TypeEndsWith   Type = "ends_with"
	TypeNotZero    Type = "not-zero"
	TypeBefore     Type = "before"
	TypeAfter      Type = "after"
	TypeDate       Type = "date"
	TypeEq         Type = "eq"
	TypeNeq        Type = "neq"
	TypeLt         Type = "lt"
	TypeLte        Type = "lte"
	TypeGt         Type = "gt"
	TypeGte        Type = "gte"
	TypeRequiredIf Type = "required_if"
)

// SupportedDirectives is the registry of all known validation directives.
var SupportedDirectives = map[Type]Info{
	TypeRequired: {
		Description:  "Field must be provided and non-empty",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "time.Time", "time.Duration", "any"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate required",
	},
	TypeOptional: {
		Description:  "Field is optional; zero values are allowed",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "time.Time", "time.Duration", "struct", "any"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate optional",
	},
	TypeMin: {
		Description:  "Minimum value for numbers, minimum length for strings/slices, minimum duration for time.Duration",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate min:18",
	},
	TypeMax: {
		Description:  "Maximum value for numbers, maximum length for strings/slices, maximum duration for time.Duration",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Duration"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate max:99",
	},
	TypeLen: {
		Description:  "Exact length for strings or exact size for slices",
		AllowedTypes: []string{"string", "slice"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate len:10",
	},
	TypeMinMax: {
		Description:  "Value range constraint — deprecated, use min and max separately",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "time.Duration"},
		ParamCount:   ParamTwo,
		Example:      "@evl:validate min-max:3,50",
		Deprecated:   true,
	},
	TypePattern: {
		Description:  "Validate string against a regular expression or predefined pattern",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate pattern:email",
		PredefinedPatterns: map[string]string{
			"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			"phone": `^\+?[0-9]{8,15}$`,
			"uuid":  `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
		},
	},
	TypeEnum: {
		Description:  "Value must be one of the specified allowed values",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64"},
		ParamCount:   ParamUnbounded,
		Example:      "@evl:validate enum:active,inactive,pending",
	},
	TypeURL: {
		Description:  "Valid URL with any scheme (ftp, http, https, etc.)",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate url",
	},
	TypeHTTPURL: {
		Description:  "Valid HTTP or HTTPS URL only",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate http_url",
	},
	TypeDSN: {
		Description:  "Valid database connection string (DSN format)",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate dsn",
	},
	TypeContains: {
		Description:  "String must contain the specified substring",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate contains:admin",
	},
	TypeStartsWith: {
		Description:  "String must start with the specified prefix",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate starts_with:https://",
	},
	TypeEndsWith: {
		Description:  "String must end with the specified suffix",
		Example:      "@evl:validate ends_with:.go",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamOne,
	},
	TypeNotZero: {
		Description:  "Value must not be zero (non-empty for slices, non-zero time for time.Time)",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64", "slice", "time.Time", "time.Duration"},
		ParamCount:   ParamNone,
		Example:      "@evl:validate not-zero",
	},
	TypeBefore: {
		Description:  "For time.Time: value must be before the specified date",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate before:2024-01-01",
	},
	TypeAfter: {
		Description:  "For time.Time: value must be after the specified date",
		AllowedTypes: []string{"time.Time"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate after:2020-01-01",
	},
	TypeDate: {
		Description:  "Validate that string is a valid date in one of the specified formats",
		AllowedTypes: []string{"string"},
		ParamCount:   ParamUnbounded,
		Example:      "@evl:validate date:RFC3339,2006-01-02",
	},
	TypeEq: {
		Description:  "Value must be equal to the specified constant",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate eq:10",
	},
	TypeNeq: {
		Description:  "Value must not be equal to the specified constant",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate neq:0",
	},
	TypeLt: {
		Description:  "Value must be less than the specified constant",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate lt:100",
	},
	TypeLte: {
		Description:  "Value must be less than or equal to the specified constant",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate lte:100",
	},
	TypeGt: {
		Description:  "Value must be greater than the specified constant",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate gt:0",
	},
	TypeGte: {
		Description:  "Value must be greater than or equal to the specified constant",
		AllowedTypes: []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"},
		ParamCount:   ParamOne,
		Example:      "@evl:validate gte:18",
	},
	TypeRequiredIf: {
		Description:  "Field is required only if another field has the specified value",
		AllowedTypes: []string{"string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "slice", "pointer", "any"},
		ParamCount:   ParamTwo, // format: "field:value"
		Example:      "@evl:validate required_if:Status active",
	},
}

// GetInfo returns metadata for a directive type, or zero Info if unknown
func GetInfo(t Type) Info {
	return SupportedDirectives[t]
}
