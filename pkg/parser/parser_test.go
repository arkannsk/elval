package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSampleFile(t *testing.T) {
	p := NewParser()
	testFile := filepath.Join("testdata", "sample", "sample.go")

	result, err := p.ParseFile(testFile)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "sample", result.Package)
	assert.Len(t, result.Structs, 4) // User, Post, Product, Config

	// Проверяем что нет ошибок валидации
	errors := result.ValidateDirectives()
	assert.Empty(t, errors)
}

func TestParseInvalidFile(t *testing.T) {
	p := NewParser()
	testFile := filepath.Join("testdata", "invalid", "invalid.go")

	result, err := p.ParseFile(testFile)
	require.NoError(t, err)
	require.NotNil(t, result)

	errors := result.ValidateDirectives()

	// Просто проверяем что есть ошибки
	assert.NotEmpty(t, errors, "должны быть ошибки валидации")

	// Выводим все ошибки для отладки
	for i, err := range errors {
		t.Logf("Ошибка %d: %v", i, err)
	}

	// Проверяем что есть хотя бы 4 ошибки
	assert.GreaterOrEqual(t, len(errors), 4, "должно быть минимум 4 ошибки")
}

func TestParseTypesFile(t *testing.T) {
	p := NewParser()
	testFile := filepath.Join("testdata", "types", "types.go")

	result, err := p.ParseFile(testFile)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "types", result.Package)
	assert.Len(t, result.Structs, 2) // AllTypes, User

	// Находим структуру AllTypes
	var allTypes *Struct
	for i := range result.Structs {
		if result.Structs[i].Name == "AllTypes" {
			allTypes = &result.Structs[i]
			break
		}
	}
	require.NotNil(t, allTypes)

	// Проверяем time.Time поле
	createdAt := allTypes.Fields[10]
	assert.Equal(t, "CreatedAt", createdAt.Name)
	assert.Equal(t, "time.Time", createdAt.Type.Name)
	assert.False(t, createdAt.Type.IsSlice)
	assert.False(t, createdAt.Type.IsPointer)
	assert.Len(t, createdAt.Directives, 3) // required, after, before

	// Проверяем директивы для time.Time
	assert.Equal(t, "required", createdAt.Directives[0].Type)
	assert.Equal(t, "after", createdAt.Directives[1].Type)
	assert.Equal(t, []string{"2020-01-01"}, createdAt.Directives[1].Params)
	assert.Equal(t, "before", createdAt.Directives[2].Type)
	assert.Equal(t, []string{"2030-12-31"}, createdAt.Directives[2].Params)

	// Проверяем time.Duration поле
	timeout := allTypes.Fields[12]
	assert.Equal(t, "Timeout", timeout.Name)
	assert.Equal(t, "time.Duration", timeout.Type.Name)
	assert.Len(t, timeout.Directives, 3) // required, min, max

	assert.Equal(t, "required", timeout.Directives[0].Type)
	assert.Equal(t, "min", timeout.Directives[1].Type)
	assert.Equal(t, []string{"1s"}, timeout.Directives[1].Params)
	assert.Equal(t, "max", timeout.Directives[2].Type)
	assert.Equal(t, []string{"24h"}, timeout.Directives[2].Params)

	// Проверяем что нет ошибок валидации
	errors := result.ValidateDirectives()
	for _, err := range errors {
		t.Logf("Validation error: %v", err)
	}
	assert.Empty(t, errors)
}

func TestParseFileNotFound(t *testing.T) {
	p := NewParser()
	testFile := filepath.Join("testdata", "not_exists", "not_exists.go")

	result, err := p.ParseFile(testFile)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка парсинга файла")
}
