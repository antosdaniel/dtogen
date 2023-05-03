package generator

import (
	"go/ast"
	"strings"
)

func determineUsedImports(expr ast.Expr, imports Imports) Imports {
	result := Imports{}
	for _, ti := range getTypeImports(expr) {
		for _, i := range imports {
			if i.UsedName() != ti {
				continue
			}
			result = append(result, i)
		}
	}
	return result
}

// getSelectors Retrieves all imports for given type
func getTypeImports(expr ast.Expr) []string {
	switch expr.(type) {
	case *ast.StarExpr:
		return getTypeImports(expr.(*ast.StarExpr).X)
	case *ast.SelectorExpr:
		// Can SelectorExpr.X be something else than Ident?
		ident, ok := expr.(*ast.SelectorExpr).X.(*ast.Ident)
		if !ok {
			return nil
		}
		return []string{ident.Name}
	case *ast.ArrayType:
		return getTypeImports(expr.(*ast.ArrayType).Elt)
	case *ast.MapType:
		m := expr.(*ast.MapType)
		return append(getTypeImports(m.Key), getTypeImports(m.Value)...)
	default:
		return nil
	}
}

type Imports []Import

type Import struct {
	alias string
	path  string
	name  string
}

func NewImport(i *ast.ImportSpec) Import {
	result := Import{
		path: strings.Trim(i.Path.Value, "\""),
	}

	if i.Name != nil {
		result.alias = i.Name.Name
	}

	result.name = getPackageName(result.path)

	return result
}

func getPackageName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (is Imports) Has(_import Import) bool {
	for _, i := range is {
		if i.path == _import.path && i.alias == _import.alias {
			return true
		}
	}

	return false
}

func (i Import) Alias() string {
	return i.alias
}

func (i Import) Path() string {
	return i.path
}

func (i Import) Name() string {
	return i.name
}

// UsedName Prefix used to access import declaration in file.
func (i Import) UsedName() string {
	if i.alias != "" {
		return i.alias
	}
	return i.name
}
