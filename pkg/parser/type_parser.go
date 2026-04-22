package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

// TypeParser отвечает за преобразование AST-выражений в FieldType
type TypeParser struct {
	typeAliases map[string]string
	currentPkg  string
	verbose     bool
}

// NewTypeParser создаёт новый парсер типов
func NewTypeParser(aliases map[string]string, pkg string, verbose bool) *TypeParser {
	return &TypeParser{
		typeAliases: aliases,
		currentPkg:  pkg,
		verbose:     verbose,
	}
}

// ParseExpr преобразует AST-выражение в FieldType
func (tp *TypeParser) ParseExpr(expr ast.Expr, allStructs map[string]*Struct) FieldType {
	var ft FieldType

	switch t := expr.(type) {
	case *ast.Ident:
		ft = tp.parseIdent(t, allStructs)

	case *ast.SelectorExpr:
		ft = tp.parseSelector(t)

	case *ast.StarExpr:
		ft = tp.parsePointer(t, allStructs)

	case *ast.ArrayType:
		ft = tp.parseArray(t, allStructs)

	case *ast.IndexExpr:
		ft = tp.parseGenericSingle(t, allStructs)

	case *ast.IndexListExpr:
		ft = tp.parseGenericMulti(t, allStructs)

	case *ast.MapType:
		ft = tp.parseMap(t)

	default:
		ft.Name = exprToString(t)
		ft.IsCustom = true
	}

	return ft
}

func (tp *TypeParser) parseIdent(t *ast.Ident, allStructs map[string]*Struct) FieldType {
	var ft FieldType
	ft.Name = t.Name

	// 1. Проверяем, структура ли это
	if _, ok := allStructs[t.Name]; ok {
		ft.IsStruct = true
		return ft
	}

	// 2. Проверяем алиасы
	if baseType, ok := tp.typeAliases[t.Name]; ok {
		ft.BaseType = baseType
		var baseExpr ast.Expr
		if idx := strings.Index(baseType, "."); idx > 0 {
			baseExpr = &ast.SelectorExpr{
				X:   &ast.Ident{Name: baseType[:idx]},
				Sel: &ast.Ident{Name: baseType[idx+1:]},
			}
		} else {
			baseExpr = &ast.Ident{Name: baseType}
		}
		baseFt := tp.ParseExpr(baseExpr, allStructs)
		ft.IsStruct = baseFt.IsStruct
		ft.IsCustom = baseFt.IsCustom
		return ft
	}

	// 3. Обычный тип
	ft.IsCustom = !isBuiltin(ft.Name)
	return ft
}

func (tp *TypeParser) parseSelector(t *ast.SelectorExpr) FieldType {
	name := exprToString(t)
	pkg := exprToString(t.X)

	// Проверяем, не является ли это известным примитивом (time.Time, time.Duration)
	isPrimitive := name == "time.Time" || name == "time.Duration"
	ft := FieldType{
		Name:        name,
		Package:     pkg,
		PackagePath: tp.currentPkg, // Для внешних типов путь текущего пакета, но для rewrite это не критично
		IsCustom:    true,
	}

	if !isPrimitive {
		ft.IsStruct = true
	}

	return ft
}

func (tp *TypeParser) parsePointer(t *ast.StarExpr, allStructs map[string]*Struct) FieldType {
	inner := tp.ParseExpr(t.X, allStructs)
	return FieldType{
		Name:        "*" + inner.Name,
		IsPointer:   true,
		IsStruct:    inner.IsStruct,
		IsCustom:    inner.IsCustom,
		IsGeneric:   inner.IsGeneric,
		GenericBase: inner.GenericBase,
		GenericArgs: inner.GenericArgs,
	}
}

func (tp *TypeParser) parseArray(t *ast.ArrayType, allStructs map[string]*Struct) FieldType {
	inner := tp.ParseExpr(t.Elt, allStructs)
	if t.Len == nil { // слайс []T
		return FieldType{
			Name:        "[]" + inner.Name,
			IsSlice:     true,
			IsStruct:    inner.IsStruct,
			IsCustom:    inner.IsCustom,
			IsGeneric:   inner.IsGeneric,
			GenericBase: inner.GenericBase,
			GenericArgs: inner.GenericArgs,
		}
	}
	lenStr := exprToString(t.Len) // например, "2" или "3"

	return FieldType{
		Name:        fmt.Sprintf("[%s]%s", lenStr, inner.Name), // "[2]float64"
		IsSlice:     true,
		IsStruct:    inner.IsStruct,
		IsCustom:    inner.IsCustom,
		IsGeneric:   inner.IsGeneric,
		GenericBase: inner.GenericBase,
		GenericArgs: inner.GenericArgs,
		BaseType:    lenStr,
	}
}

func (tp *TypeParser) parseGenericSingle(t *ast.IndexExpr, allStructs map[string]*Struct) FieldType {
	base := exprToString(t.X)
	inner := tp.ParseExpr(t.Index, allStructs)
	return FieldType{
		Name:        fmt.Sprintf("%s[%s]", base, inner.Name),
		IsGeneric:   true,
		GenericBase: base,
		GenericArgs: []FieldType{inner},
		IsCustom:    true,
	}
}

func (tp *TypeParser) parseGenericMulti(t *ast.IndexListExpr, allStructs map[string]*Struct) FieldType {
	base := exprToString(t.X)
	var args []FieldType
	var argNames []string
	for _, idx := range t.Indices {
		arg := tp.ParseExpr(idx, allStructs)
		args = append(args, arg)
		argNames = append(argNames, arg.Name)
	}
	return FieldType{
		Name:        fmt.Sprintf("%s[%s]", base, strings.Join(argNames, ", ")),
		IsGeneric:   true,
		GenericBase: base,
		GenericArgs: args,
		IsCustom:    true,
	}
}

func (tp *TypeParser) parseMap(t *ast.MapType) FieldType {
	return FieldType{
		Name:     fmt.Sprintf("map[%s]%s", exprToString(t.Key), exprToString(t.Value)),
		IsCustom: true,
	}
}

// exprToString преобразует AST-выражение в строку
func exprToString(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprToString(v.X) + "." + v.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(v.X)
	case *ast.ArrayType:
		if v.Len == nil {
			return "[]" + exprToString(v.Elt)
		}
		// Для массивов [N]T рекурсивно обрабатываем длину и элемент
		lenStr := exprToString(v.Len)
		return fmt.Sprintf("[%s]%s", lenStr, exprToString(v.Elt))
	case *ast.IndexExpr:
		return exprToString(v.X) + "[" + exprToString(v.Index) + "]"
	case *ast.IndexListExpr:
		var parts []string
		for _, idx := range v.Indices {
			parts = append(parts, exprToString(idx))
		}
		return exprToString(v.X) + "[" + strings.Join(parts, ", ") + "]"
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(v.Key), exprToString(v.Value))
	case *ast.ParenExpr:
		return "(" + exprToString(v.X) + ")"
	case *ast.BasicLit:
		return v.Value

	default:
		return "unknown"
	}
}
