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

		var usedImports Imports
		if desired.OverrideTypeTo == "" {
			usedImports = determineUsedImports(f.Type(), imports)
		}

		result = append(result, Field{
			name:         f.Name(),
			renameTo:     desired.RenameTo,
			_type:        f.Type(),
			overrideType: desired.OverrideTypeTo,
			imports:      usedImports,
		})
	}

	return result
}

type Fields []Field

type Field struct {
	name         string
	renameTo     string
	_type        ast.Expr
	overrideType string
	imports      Imports
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

func (f Field) Type() ast.Expr {
	return f._type
}

func (f Field) OverrideTypeTo() string {
	return f.overrideType
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
