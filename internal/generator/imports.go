package generator

import (
	"go/ast"
	"strings"
)

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

	parts := strings.Split(result.path, "/")
	result.name = parts[len(parts)-1]

	return result
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
