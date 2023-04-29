package internal

import (
	"github.com/antosdaniel/dtogen/internal/parser"
	"github.com/antosdaniel/dtogen/internal/writer"
	"go/ast"
	"strings"
)

type Generated interface {
	String() string
}

func Generate(input Input) (Generated, error) {
	p, err := parser.New(input.PathToSource)
	if err != nil {
		return nil, err
	}
	parsed, err := p.GetStruct(input.TypeName)
	if err != nil {
		return nil, err
	}

	// TODO: error on not found names
	fieldsToWrite := writer.Fields{}
	for _, f := range prepareFields(input.IncludeAllParsedFields, input.Fields, parsed.Fields) {
		fieldsToWrite = append(fieldsToWrite, writer.Field{
			Name: f.desiredName(),
			Type: f.Type,
		})
	}
	w := writer.New()
	w.WritePackage(input.OutputPackage)
	w.WriteEmptyLine()
	// TODO: imports
	w.WriteStruct(writer.Struct{
		Name:   input.desiredTypeName(),
		Fields: fieldsToWrite,
	})

	return w, nil
}

func prepareFields(includeAllParsedFields bool, desiredFields FieldsInput, parsedFields parser.Fields) fields {
	result := fields{}
	for _, f := range parsedFields {
		desired, found := desiredFields.findByOriginalName(f.Name)
		if !found && !includeAllParsedFields {
			continue
		}

		result = append(result, field{
			Name:     f.Name,
			RenameTo: desired.RenameTo,
			Type:     f.Type,
		})
	}

	return result
}

type fields []field

type field struct {
	Name     string
	RenameTo string
	Type     ast.Expr
}

func (f field) desiredName() string {
	if f.RenameTo != "" {
		return f.RenameTo
	}

	if !ast.IsExported(f.Name) {
		return strings.ToUpper(f.Name[:1]) + f.Name[1:]
	}

	return f.Name
}
