package generator

import (
	"go/ast"
	"strings"
)

// TODO: Refactor `possibleMappers` and `helpers`. They are kind of a hack right now.
//  Helpers should be preserved somewhere else. This function should only try to find handwritten mappers names.

func prepareMapper(
	destinationTypeName string,
	parsed ParsedStruct,
	fields Fields,
	possibleMappers ParsedFunctions,
) Mapper {
	result := Mapper{
		SourceTypeName:      parsed.Name,
		SourceImportName:    parsed.Package.name,
		DestinationTypeName: destinationTypeName,
	}

	for _, f := range fields {
		mapping, helper := mappingForField(parsed, f, possibleMappers)
		result.Mappings = append(result.Mappings, mapping)
		if helper != nil {
			result.Helpers = append(result.Helpers, *helper)
		}
	}

	return result
}

func mappingForField(parsed ParsedStruct, field Field, possibleMappers ParsedFunctions) (Mapping, *ParsedFunction) {
	mapping := Mapping{destinationFieldName: field.DesiredName()}

	for _, m := range possibleMappers {
		f := HandwrittenMapperName(mapping.destinationFieldName)
		if strings.ToLower(m.name) == strings.ToLower(f) {
			mapping.sourceFunction = m.name
			mapping.sourceField = field.OriginalName() // TODO: this is a hack
			tmp := m
			return mapping, &tmp
		}
	}

	// If field we are mapping from is exported, we can just assign it.
	if ast.IsExported(field.OriginalName()) {
		mapping.sourceField = field.OriginalName()
		return mapping, nil
	}

	// Field is not exported, we need to look for `Field()` or `GetField()`.
	// TODO: Skip methods with arguments
	method, found := parsed.Methods.findExportedMethod(
		field.OriginalName(),
		"get"+field.OriginalName(),
	)
	if found {
		mapping.sourceMethod = method.name
		return mapping, nil
	}

	return mapping, nil
}

type Mapper struct {
	// SourceTypeName Name of a struct we are mapping from.
	SourceTypeName string
	// SourceImportName Name of an import of struct we are mapping from.
	SourceImportName string

	// DestinationTypeName Name of a struct we are mapping to.
	DestinationTypeName string

	Mappings Mappings

	Helpers ParsedFunctions
}

type Mappings []Mapping

type Mapping struct {
	destinationFieldName string

	// TODO: refactor sources' types
	sourceField    string
	sourceMethod   string
	sourceFunction string
}

func (m Mapping) Destination() string {
	return m.destinationFieldName
}

func (m Mapping) Source() string {
	if m.IsField() {
		return m.sourceField
	}
	if m.IsFunction() {
		return m.sourceFunction
	}

	return m.sourceMethod
}

func (m Mapping) Field() string {
	return m.sourceField
}

func (m Mapping) IsField() bool {
	return m.sourceField != "" && !m.IsFunction() // TODO: this is a hack
}

func (m Mapping) IsMethod() bool {
	return m.sourceMethod != ""
}

func (m Mapping) IsFunction() bool {
	return m.sourceFunction != ""
}

func HandwrittenMapperName(destinationFieldName string) string {
	return "new" + destinationFieldName
}
