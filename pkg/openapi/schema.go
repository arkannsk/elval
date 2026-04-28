package openapi

// Discriminator описывает полиморфную диспетчеризацию
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// Schema представляет объект схемы OpenAPI 3.0/3.1
type Schema struct {
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Example     any    `json:"example,omitempty"`

	// Для объектов
	Properties map[string]*Schema `json:"properties,omitempty"` // 👈 ИЗМЕНЕНО: *Schema
	Required   []string           `json:"required,omitempty"`

	// Для массивов
	Items *Schema `json:"items,omitempty"` // 👈 Уже было правильно

	// Перечисления и ограничения
	Enum             []any    `json:"enum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	MinLength        *int64   `json:"minLength,omitempty"`
	MaxLength        *int64   `json:"maxLength,omitempty"`
	Pattern          string   `json:"pattern,omitempty"`

	// Полиморфизм
	Discriminator *Discriminator `json:"discriminator,omitempty"`
	OneOf         []*Schema      `json:"oneOf,omitempty"`
	AllOf         []*Schema      `json:"allOf,omitempty"`
	AnyOf         []*Schema      `json:"anyOf,omitempty"`

	// Ссылка на компонент
	Ref string `json:"$ref,omitempty"`

	// Дополнительные флаги
	Nullable  bool `json:"nullable,omitempty"`
	ReadOnly  bool `json:"readOnly,omitempty"`
	WriteOnly bool `json:"writeOnly,omitempty"`
	Default   any  `json:"default,omitempty"`
}

// NewSchema создает новую схему с инициализированными коллекциями
func NewSchema() *Schema {
	return &Schema{
		Properties: make(map[string]*Schema), // 👈 Инициализация с указателями
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
	Schema      *Schema   `json:"schema,omitempty"` // 👈 Добавлен omitempty, так как schema может быть nil для некоторых параметров (редко, но бывает)
	Example     any       `json:"example,omitempty"`

	// Дополнительные поля из спецификации Parameter Object
	Deprecated      bool   `json:"deprecated,omitempty"`
	AllowEmptyValue bool   `json:"allowEmptyValue,omitempty"`
	Style           string `json:"style,omitempty"`
	Explode         *bool  `json:"explode,omitempty"`
}
