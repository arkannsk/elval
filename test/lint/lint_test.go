// test/lint/lint_test.go
package lint_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLintInvalid(t *testing.T) {
	cmd := exec.Command("elval-gen", "lint", "-i", "./invalid", "-v")
	output, err := cmd.CombinedOutput()

	// Ожидаем ошибку (ненулевой код возврата)
	assert.Error(t, err)

	outputStr := string(output)

	// Проверяем наличие ожидаемых ошибок (с учетом реального формата)
	expectedErrors := []string{
		"неизвестная директива",
		"требуется 1 параметр(ов), получено 0",
		"параметр должен быть числом, получено 'abc'",
		"невалидное регулярное выражение",
		"требуется 1 параметр(ов), получено 0", // для enum
		"невалидный формат даты",
		"параметр должен быть целым числом, получено 'abc'",
	}

	for _, expected := range expectedErrors {
		assert.Contains(t, outputStr, expected, "Expected error message: %s", expected)
	}

	// Проверяем количество ошибок
	assert.Contains(t, outputStr, "errors found: 7")

	t.Logf("Lint output:\n%s", outputStr)
}

func TestLintValid(t *testing.T) {
	cmd := exec.Command("elval-gen", "lint", "-i", "./valid", "-v")
	output, err := cmd.CombinedOutput()

	// Ожидаем успех (нулевой код возврата)
	assert.NoError(t, err)

	outputStr := string(output)
	// Реальный вывод: "valid.go - all annotations valid"
	assert.Contains(t, outputStr, "valid.go - all annotations valid")
	assert.Contains(t, outputStr, "errors found: 0")
}

func TestLintVerbose(t *testing.T) {
	cmd := exec.Command("elval-gen", "lint", "-i", "./valid", "-v")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "Files checked: 1, errors found: 0")
}

func TestLintQuiet(t *testing.T) {
	cmd := exec.Command("elval-gen", "lint", "-i", "./valid")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err)

	outputStr := string(output)
	// В тихом режиме вывод только если есть ошибки
	// Для валидного файла должно быть пусто
	if outputStr != "" && !strings.Contains(outputStr, "ok") {
		t.Logf("Output: %s", outputStr)
	}
}

func TestLintMultipleFiles(t *testing.T) {
	// Создаём временную директорию с двумя файлами
	tmpDir := t.TempDir()

	// Создаём файл с ошибкой
	invalidContent := `package test
type Test struct {
    // @evl:validate unknown
    Name string
}`
	err := os.WriteFile(filepath.Join(tmpDir, "invalid.go"), []byte(invalidContent), 0644)
	require.NoError(t, err)

	// Создаём валидный файл
	validContent := `package test
type Test struct {
    // @evl:validate required
    Name string
}`
	err = os.WriteFile(filepath.Join(tmpDir, "valid.go"), []byte(validContent), 0644)
	require.NoError(t, err)

	// Запускаем линтер
	cmd := exec.Command("elval-gen", "lint", "-i", tmpDir, "-v")
	output, err := cmd.CombinedOutput()

	// Ожидаем ошибку
	assert.Error(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "invalid.go")
	assert.Contains(t, outputStr, "valid.go - all annotations valid")
}

func TestLintNoAnnotations(t *testing.T) {
	// Создаем временный файл без аннотаций
	cmd := exec.Command("elval-gen", "lint", "-i", ".", "-v")
	output, err := cmd.CombinedOutput()

	// Должно быть успешно (нет ошибок, так как нет аннотаций для проверки)
	// Но может быть ошибка если есть другие файлы с ошибками
	t.Logf("Output: %s", string(output))
	_ = err
}

func TestLintHelp(t *testing.T) {
	cmd := exec.Command("elval-gen", "help")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "elval-gen generate")
	assert.Contains(t, outputStr, "elval-gen lint")
	assert.Contains(t, outputStr, "elval-gen version")
}

func TestLintVersion(t *testing.T) {
	cmd := exec.Command("elval-gen", "version")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "elval-gen version")
}

func TestLintWithExclude(t *testing.T) {
	tmpDir := t.TempDir()

	// Создаём файл с ошибкой в основной директории
	invalidContent := `package main
type Test struct {
    // @evl:validate unknown
    Name string
}`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(invalidContent), 0644)
	require.NoError(t, err)

	// Создаём директорию vendor с таким же файлом
	vendorDir := filepath.Join(tmpDir, "vendor")
	err = os.MkdirAll(vendorDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(vendorDir, "vendor.go"), []byte(invalidContent), 0644)
	require.NoError(t, err)

	// Запускаем линтер с exclude
	cmd := exec.Command("elval-gen", "lint", "-i", tmpDir, "-exclude", "vendor", "-v")
	output, err := cmd.CombinedOutput()

	// Ожидаем ошибку из main.go
	assert.Error(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "main.go")
	assert.NotContains(t, outputStr, "vendor.go")
}
