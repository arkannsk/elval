package oa

type Schema struct {
	Type             string            `json:"type,omitempty"`
	Format           string            `json:"format,omitempty"`
	Properties       map[string]Schema `json:"properties,omitempty"`
	Required         []string          `json:"required,omitempty"`
	MinLength        *int64            `json:"minLength,omitempty"`
	MaxLength        *int64            `json:"maxLength,omitempty"`
	Minimum          *float64          `json:"minimum,omitempty"`
	Maximum          *float64          `json:"maximum,omitempty"`
	ExclusiveMinimum bool              `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool              `json:"exclusiveMaximum,omitempty"`
	MinItems         *int64            `json:"minItems,omitempty"`
	MaxItems         *int64            `json:"maxItems,omitempty"`
	Pattern          string            `json:"pattern,omitempty"`
	Enum             []interface{}     `json:"enum,omitempty"`
	Title            string            `json:"title,omitempty"`
	Description      string            `json:"description,omitempty"`
	Example          interface{}       `json:"example,omitempty"`
	Items            *Schema           `json:"items,omitempty"`
	Ref              string            `json:"$ref,omitempty"`
}

// NewSchema creates a new Schema
func NewSchema() *Schema {
	return &Schema{
		Properties: make(map[string]Schema),
		Required:   []string{},
	}
}
