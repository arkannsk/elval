package oa

// Discriminator описывает полиморфную диспетчеризацию
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// Schema расширяем для поддержки OpenAPI 3.0 polymorphism
type Schema struct {
	Type             string            `json:"type,omitempty"`
	Format           string            `json:"format,omitempty"`
	Title            string            `json:"title,omitempty"`
	Description      string            `json:"description,omitempty"`
	Example          any               `json:"example,omitempty"`
	Properties       map[string]Schema `json:"properties,omitempty"`
	Required         []string          `json:"required,omitempty"`
	Enum             []any             `json:"enum,omitempty"`
	Minimum          *float64          `json:"minimum,omitempty"`
	Maximum          *float64          `json:"maximum,omitempty"`
	ExclusiveMinimum bool              `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool              `json:"exclusiveMaximum,omitempty"`
	MinLength        *int64            `json:"minLength,omitempty"`
	MaxLength        *int64            `json:"maxLength,omitempty"`
	Pattern          string            `json:"pattern,omitempty"`
	Items            *Schema           `json:"items,omitempty"`

	Discriminator *Discriminator `json:"discriminator,omitempty"`
	OneOf         []Schema       `json:"oneOf,omitempty"`
	AllOf         []Schema       `json:"allOf,omitempty"`
	Ref           string         `json:"$ref,omitempty"`
}

// NewSchema creates a new Schema
func NewSchema() *Schema {
	return &Schema{
		Properties: make(map[string]Schema),
		Required:   []string{},
	}
}

// ParamType тип параметра
type ParamType string

const (
	ParamPath   ParamType = "path"
	ParamQuery  ParamType = "query"
	ParamHeader ParamType = "header"
	ParamCookie ParamType = "cookie"
)

// Parameter описание параметра для OpenAPI
type Parameter struct {
	Name        string    `json:"name"`
	In          ParamType `json:"in"`
	Description string    `json:"description,omitempty"`
	Required    bool      `json:"required,omitempty"`
	Schema      *Schema   `json:"schema"`
	Example     any       `json:"example,omitempty"`
}
