package generator

import (
	"go/ast"
	"strings"
)

func prepareFields(
	includeAllParsedFields bool,
	desiredFields FieldsInput,
	parsedFields ParsedFields,
	imports Imports,
) Fields {
	result := Fields{}
	for _, f := range parsedFields {
		desired, found := desiredFields.findByOriginalName(f.Name())
		if !found && !includeAllParsedFields {
			continue
		}

		result = append(result, Field{
			name:     f.Name(),
			renameTo: desired.RenameTo,
			_type:    f.Type(),
			imports:  determineUsedImports(f.Type(), imports),
		})
	}

	return result
}

func determineUsedImports(_type ast.Expr, imports Imports) Imports {
	expr := _type
	// If it's a pointer, get underlying type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}

	sel, ok := expr.(*ast.SelectorExpr)
	// TODO: revisit with lists/maps
	if !ok {
		return nil
	}

	result := Imports{}
	for _, i := range imports {
		// Can it be something else?
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}
		if i.UsedName() != ident.Name {
			continue
		}
		result = append(result, i)
	}
	return result
}

type Fields []Field
type Field struct {
	name     string
	renameTo string
	_type    ast.Expr
	imports  Imports
}

func (f Field) Type() ast.Expr {
	return f._type
}

func (f Field) DesiredName() string {
	if f.renameTo != "" {
		return f.renameTo
	}

	if !ast.IsExported(f.name) {
		return strings.ToUpper(f.name[:1]) + f.name[1:]
	}

	return f.name
}

func (fs Fields) RequiredImports() Imports {
	result := Imports{}
	for _, f := range fs {
		for _, i := range f.imports {
			if !result.Has(i) {
				result = append(result, i)
			}
		}
	}

	return result
}
