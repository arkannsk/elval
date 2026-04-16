package parser

// FieldType представляет тип поля с поддержкой слайсов и указателей
type FieldType struct {
	Name      string
	IsSlice   bool
	IsPointer bool
	IsStruct  bool
}

// String возвращает строковое представление типа
func (ft FieldType) String() string {
	if ft.IsSlice {
		return "[]" + ft.Name
	}
	if ft.IsPointer {
		return "*" + ft.Name
	}
	return ft.Name
}

// Directive представляет одну аннотацию валидации
type Directive struct {
	Type   string   // required, min, max, len, pattern, not-zero, optional
	Params []string // параметры директивы
	Raw    string   // исходный текст
}

type OaAnnotation struct {
	Type  string // title, description, example, format, etc.
	Value string
}

// Field представляет поле структуры с аннотациями
type Field struct {
	Name          string      // имя поля
	Type          FieldType   // тип поля
	Directives    []Directive // список директив валидации
	Decorators    []Decorator // декораторы
	Line          int         // номер строки в файле (для ошибок)
	OaAnnotations []OaAnnotation
}

// Struct представляет структуру с полями для валидации
type Struct struct {
	Name   string  // имя структуры
	Fields []Field // поля с аннотациями
	File   string  // путь к файлу
}

// ParseResult результат парсинга файла
type ParseResult struct {
	Package string   // имя пакета
	Structs []Struct // найденные структуры
	Errors  []error  // ошибки парсинга
}
