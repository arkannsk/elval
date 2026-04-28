package parser

import (
	"go/parser"
	"go/token"
	"testing"

	ann "github.com/arkannsk/elval/pkg/parser/annotations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helperParseString помогает парсить код из строки для тестов
func helperParseString(t *testing.T, code string) *ParseResult {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	assert.NoError(t, err)

	p := NewParser(false) // verbose=false для чистоты логов

	modInfo := &ModuleInfo{
		Module:      "github.com/test/module",
		ModuleRoot:  "/tmp/test",
		PackagePath: "test/pkg",
		Package:     node.Name.Name,
	}

	result := &ParseResult{
		Package: node.Name.Name,
		Structs: make([]Struct, 0),
		Errors:  make([]error, 0),
	}

	allStructs, typeAliases := p.parseStructsFirstPass(node, "test.go", modInfo)
	p.typeParser = NewTypeParser(typeAliases, modInfo.Package, false)

	p.parseFieldsSecondPass(node, allStructs, result, "test.go", modInfo)

	return result
}

func TestParser_ParseSimpleStruct(t *testing.T) {
	code := `
package main

type User struct {
	// @evl:validate required
	ID   int
	
	// @oa:description Имя пользователя
	Name string
}
`

	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs, "Структура User должна быть распаршена")

	if len(result.Structs) == 0 {
		t.Skip("Структура не найдена")
	}

	user := result.Structs[0]
	assert.Equal(t, "User", user.Name)
	assert.Len(t, user.Fields, 2)

	// Проверка поля ID
	idField := user.Fields[0]
	assert.Equal(t, "ID", idField.Name)
	assert.Equal(t, "int", idField.Type.Name)

	// Проверяем директивы
	assert.Len(t, idField.Directives, 1, "Ожидаем одну директиву 'required'")
	assert.Equal(t, "required", idField.Directives[0].Type)

	// Проверка поля Name
	nameField := user.Fields[1]
	assert.Equal(t, "Name", nameField.Name)
	assert.Equal(t, "string", nameField.Type.Name)

	// Проверяем OA аннотации
	assert.Len(t, nameField.OaAnnotations, 1, "Ожидаем одну аннотацию 'description'")
	assert.Equal(t, "description", nameField.OaAnnotations[0].Type)
	assert.Equal(t, "Имя пользователя", nameField.OaAnnotations[0].Value)
}

func TestParser_ParseIgnoredStruct(t *testing.T) {
	// Аннотация ignore на уровне структуры
	code := `
package main

// @oa:ignore
type InternalData struct {
	Secret string
}
`

	result := helperParseString(t, code)

	// Структура должна быть пропущена
	assert.Empty(t, result.Structs)
}

func TestParser_ParseIgnoredField(t *testing.T) {
	// Аннотация ignore на поле
	code := `
package main

type User struct {
	PublicID string
	// @oa:ignore
	PrivateKey string
}
`

	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs)
	user := result.Structs[0]

	// Должно остаться только одно поле (PublicID), PrivateKey игнорируется
	assert.Len(t, user.Fields, 1)
	assert.Equal(t, "PublicID", user.Fields[0].Name)
}

func TestParser_ParseDiscriminator(t *testing.T) {
	code := `
package main

// @oa:discriminator type
// @oa:oneOf Dog,Cat
type Animal struct {
	Type string
}
`
	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs)
	animal := result.Structs[0]

	require.NotNil(t, animal.Discriminator, "Discriminator не должен быть nil")
	assert.Equal(t, "type", animal.Discriminator.PropertyName)
	assert.Equal(t, []string{"Dog", "Cat"}, animal.OaOneOf)
}

func TestParser_ParseFieldOAIn(t *testing.T) {
	code := `
package main

type GetUserRequest struct {
	// @oa:in path id
	UserID string
	
	// @oa:in query fields
	Fields []string
}
`

	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs)
	req := result.Structs[0]
	assert.Len(t, req.Fields, 2)

	// Проверка UserID
	userIDField := req.Fields[0]
	assert.Equal(t, "path", userIDField.OaIn)
	assert.Equal(t, "id", userIDField.OaParamName)

	// Проверка Fields
	fieldsField := req.Fields[1]
	assert.Equal(t, "query", fieldsField.OaIn)
	assert.Equal(t, "fields", fieldsField.OaParamName)
}

func TestParser_ParseRewriteRef(t *testing.T) {
	// Аннотация rewrite.ref на поле
	code := `
package main

type Response struct {
	// @oa:rewrite.ref #/components/schemas/User
	Data interface{}
}
`

	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs)
	resp := result.Structs[0]
	dataField := resp.Fields[0]

	assert.Equal(t, "#/components/schemas/User", dataField.OaRewriteRef)
}

func TestParser_ParseNestedStruct(t *testing.T) {
	code := `
package main

type Address struct {
	City string
}

type User struct {
	ID int
	Addr Address
}
`

	result := helperParseString(t, code)

	assert.Len(t, result.Structs, 2)

	// Находим User
	var user *Struct
	for i := range result.Structs {
		if result.Structs[i].Name == "User" {
			user = &result.Structs[i]
			break
		}
	}
	assert.NotNil(t, user)

	addrField := user.Fields[1]
	assert.Equal(t, "Address", addrField.Type.Name)
	assert.True(t, addrField.Type.IsStruct)
}

func TestParser_ParseSliceField(t *testing.T) {
	code := `
package main

type Item struct {
	Name string
}

type Order struct {
	Items []Item
}
`

	result := helperParseString(t, code)

	assert.Len(t, result.Structs, 2)

	var order *Struct
	for i := range result.Structs {
		if result.Structs[i].Name == "Order" {
			order = &result.Structs[i]
			break
		}
	}
	assert.NotNil(t, order)

	itemsField := order.Fields[0]
	assert.True(t, itemsField.Type.IsSlice)
	assert.Equal(t, "Item", itemsField.Type.GenericBase)
}

func TestParser_ParseComplexValidation(t *testing.T) {
	// Используем блок комментариев или несколько строк // для читаемости
	code := `
package main

type Product struct {
	// @evl:validate min:1
	// @evl:validate max:100
	// @evl:validate pattern:^\d+$
	Quantity int
}
`

	result := helperParseString(t, code)

	require.NotEmpty(t, result.Structs)
	product := result.Structs[0]
	qtyField := product.Fields[0]

	assert.Equal(t, "Quantity", qtyField.Name)

	// Проверяем количество директив
	require.Len(t, qtyField.Directives, 3, "Ожидаем три директивы: min, max, pattern")

	// Сортируем или ищем по типу, так как порядок может зависеть от парсера
	// Но обычно они идут в порядке появления в коде

	minDir := findDirective(qtyField.Directives, "min")
	require.NotNil(t, minDir)
	assert.Equal(t, []string{"1"}, minDir.Params)

	maxDir := findDirective(qtyField.Directives, "max")
	require.NotNil(t, maxDir)
	assert.Equal(t, []string{"100"}, maxDir.Params)

	patternDir := findDirective(qtyField.Directives, "pattern")
	require.NotNil(t, patternDir)
	assert.Equal(t, []string{"^\\d+$"}, patternDir.Params)
}

func TestParser_ParseMultipleOAAnnotations(t *testing.T) {
	// Несколько OA аннотаций на одном поле
	code := `
package main

type APIResponse struct {
	// @oa:description Список пользователей
	// @oa:format array
	Users []string
}
`

	result := helperParseString(t, code)

	assert.NotEmpty(t, result.Structs)
	resp := result.Structs[0]
	usersField := resp.Fields[0]

	assert.Equal(t, "Users", usersField.Name)
	assert.Len(t, usersField.OaAnnotations, 2)

	// Порядок может зависеть от реализации, проверим наличие типов
	types := make(map[string]bool)
	for _, an := range usersField.OaAnnotations {
		types[an.Type] = true
	}
	assert.True(t, types["description"])
	assert.True(t, types["format"])
}

func findDirective(directives []ann.Directive, name string) *ann.Directive {
	for i := range directives {
		if directives[i].Type == name {
			return &directives[i]
		}
	}
	return nil
}
