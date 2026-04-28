package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModuleInfo хранит информацию о модуле и пакете
type ModuleInfo struct {
	Module      string // "github.com/arkannsk/elval"
	Package     string // "user"
	PackagePath string // "test/integration/oa_unique_refs/user"
	ModuleRoot  string // "/home/user/elval" (absolute)
}

// resolveModuleInfo вычисляет информацию о модуле для данного файла
func resolveModuleInfo(filename string) (*ModuleInfo, error) {
	absFile, err := filepath.Abs(filename)
	if err != nil {
		absFile = filename
	}

	info := &ModuleInfo{}

	// 1. Находим go.mod и извлекаем модуль
	info.Module = getModulePath(absFile)
	if info.Module == "" {
		return nil, fmt.Errorf("go.mod not found for %s", filename)
	}

	// 2. Находим корень модуля на диске
	root, err := findModuleRoot(absFile)
	if err != nil {
		return nil, err
	}
	info.ModuleRoot = root

	// 3. Вычисляем путь пакета относительно корня модуля
	info.PackagePath = getPackagePath(absFile, root)

	return info, nil
}

// getModulePath читает модуль из go.mod в директории файла
func getModulePath(filename string) string {
	dir := filepath.Dir(filename)
	for {
		modFile := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(modFile); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module "))
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// findModuleRoot ищет директорию с go.mod, поднимаясь вверх от filename
func findModuleRoot(filename string) (string, error) {
	dir := filepath.Dir(filename)
	for {
		modFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			return filepath.Abs(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found for %s", filename)
		}
		dir = parent
	}
}

// getPackagePath вычисляет путь пакета относительно корня модуля
func getPackagePath(filename, moduleRoot string) string {
	if moduleRoot == "" {
		return filepath.Base(filepath.Dir(filename))
	}
	fileDir := filepath.Dir(filename)
	rel, err := filepath.Rel(moduleRoot, fileDir)
	if err != nil || rel == "." {
		return filepath.Base(fileDir)
	}
	return filepath.ToSlash(rel)
}
