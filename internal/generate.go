package internal

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/writer"
)

var (
	ErrNameCanNotBeEmpty error = fmt.Errorf("name can not be empty")
)

type FieldName string
type FieldType string

type Generated interface {
	String() string
}

func Generate(input Input) (Generated, error) {
	parsed, err := Parse(input.PathToSource, input.TypeName)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %q: %w", input.PathToSource, err)
	}

	// TODO: error on not found names
	fieldsToWrite := writer.Fields{}
	for _, f := range prepareFields(input.IncludeAllParsedFields, input.Fields, parsed.Fields) {
		fieldsToWrite = append(fieldsToWrite, writer.Field{
			Name: string(f.desiredName()),
			Type: string(f.Type),
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

func prepareFields(includeAllParsedFields bool, desiredFields FieldsInput, parsedFields ParsedFields) fields {
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
	Name     FieldName
	RenameTo FieldName
	Type     FieldType
}

func (f field) desiredName() FieldName {
	if f.RenameTo != "" {
		return f.RenameTo
	}

	return f.Name
}
