package openapi

import (
	"testing"

	"github.com/arkannsk/elval/pkg/oa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserOaSchema(t *testing.T) {
	user := User{}
	schema := user.OaSchema()

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 4)
	assert.ElementsMatch(t, []string{"name", "email", "age", "phone"}, getPropertyNames(schema))

	// Проверяем поле Name
	nameProp, ok := schema.Properties["name"]
	require.True(t, ok)
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, int64(3), *nameProp.MinLength)
	assert.Equal(t, int64(50), *nameProp.MaxLength)
	assert.Equal(t, "User Name", nameProp.Title)
	assert.Equal(t, "Full name of the user", nameProp.Description)
	assert.Equal(t, "John Doe", nameProp.Example)

	// Проверяем поле Email
	emailProp, ok := schema.Properties["email"]
	require.True(t, ok)
	assert.Equal(t, "string", emailProp.Type)
	assert.Equal(t, "email", emailProp.Format)
	assert.Contains(t, schema.Required, "email")

	// Проверяем поле Age
	ageProp, ok := schema.Properties["age"]
	require.True(t, ok)
	assert.Equal(t, "integer", ageProp.Type)
	assert.Equal(t, float64(18), *ageProp.Minimum)
	assert.Equal(t, float64(120), *ageProp.Maximum)

	// Проверяем поле Phone (optional)
	phoneProp, ok := schema.Properties["phone"]
	require.True(t, ok)
	assert.Equal(t, "string", phoneProp.Type)
	assert.Equal(t, "phone", phoneProp.Format)
	assert.NotContains(t, schema.Required, "phone")
}

func TestProductOaSchema(t *testing.T) {
	product := Product{}
	schema := product.OaSchema()

	require.NotNil(t, schema)

	if itemsProp, ok := schema.Properties["items"]; ok {
		require.NotNil(t, itemsProp.Items, "Items.Items не должен быть nil")
	}

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 4)

	// Проверяем поле Status
	statusProp, ok := schema.Properties["status"]
	require.True(t, ok)
	assert.Equal(t, "string", statusProp.Type)
	assert.ElementsMatch(t, []interface{}{"active", "inactive", "archived"}, statusProp.Enum)
	assert.Equal(t, "Product status", statusProp.Description)
	assert.Contains(t, schema.Required, "status")

	// Проверяем поле Price
	priceProp, ok := schema.Properties["price"]
	require.True(t, ok)
	assert.Equal(t, "number", priceProp.Type)
	assert.Equal(t, float64(0), *priceProp.Minimum)
	assert.Contains(t, schema.Required, "price")

	// Проверяем поле Quantity (optional)
	qtyProp, ok := schema.Properties["quantity"]
	require.True(t, ok)
	assert.Equal(t, "integer", qtyProp.Type)
	assert.Equal(t, float64(1), *qtyProp.Minimum)
	assert.Equal(t, float64(1000), *qtyProp.Maximum)
	assert.NotContains(t, schema.Required, "quantity")

	// Проверяем поле Code
	codeProp, ok := schema.Properties["code"]
	require.True(t, ok)
	assert.Equal(t, "string", codeProp.Type)
	assert.Equal(t, int64(10), *codeProp.MinLength)
	assert.Equal(t, int64(10), *codeProp.MaxLength)
	assert.Contains(t, schema.Required, "code")
}

func TestOrderOaSchema(t *testing.T) {
	order := Order{}
	schema := order.OaSchema()

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 3)

	// Проверяем поле ID
	idProp, ok := schema.Properties["id"]
	require.True(t, ok)
	assert.Equal(t, "string", idProp.Type)
	assert.Contains(t, schema.Required, "id")

	// Проверяем поле CreatedAt
	createdAtProp, ok := schema.Properties["createdat"]
	require.True(t, ok)
	assert.Equal(t, "string", createdAtProp.Type)
	assert.Equal(t, "date-time", createdAtProp.Format)
	assert.Contains(t, schema.Required, "createdat")

	// Проверяем поле Items (slice)
	itemsProp, ok := schema.Properties["items"]
	require.True(t, ok)
	assert.Equal(t, "array", itemsProp.Type)
	assert.NotNil(t, itemsProp.Items)
	assert.Equal(t, "string", itemsProp.Items.Type)
	assert.Equal(t, int64(1), *itemsProp.MinLength)
	assert.Equal(t, int64(100), *itemsProp.MaxLength)
	assert.Contains(t, schema.Required, "items")
}

func getPropertyNames(schema *oa.Schema) []string {
	names := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		names = append(names, name)
	}
	return names
}
