package generator

import (
	"go/ast"
)

func prepareMapper(
	parsed ParsedStruct,
	fields Fields,
	destinationTypeName string,
) Mapper {
	result := Mapper{
		SourceTypeName:      parsed.Name,
		SourceImportName:    parsed.Package.name,
		DestinationTypeName: destinationTypeName,
	}

	for _, f := range fields {
		result.Mappings = append(result.Mappings, mappingForField(parsed, f))
	}

	return result
}

func mappingForField(parsed ParsedStruct, field Field) Mapping {
	mapping := Mapping{destinationFieldName: field.DesiredName()}
	// If field we are mapping from is exported, we can just assign it.
	if ast.IsExported(field.OriginalName()) {
		mapping.sourceField = field.OriginalName()
		return mapping
	}

	// Field is not exported, we need to look for `Field()` or `GetField()`.
	// TODO: Skip methods with arguments
	method, found := parsed.Methods.findExportedMethod(
		field.OriginalName(),
		"get"+field.OriginalName(),
	)
	if found {
		mapping.sourceMethod = method.name
		return mapping
	}

	return mapping
}

type Mapper struct {
	// SourceTypeName Name of a struct we are mapping from.
	SourceTypeName string
	// SourceImportName Name of an import of struct we are mapping from.
	SourceImportName string

	// DestinationTypeName Name of a struct we are mapping to.
	DestinationTypeName string

	Mappings Mappings
}

type Mappings []Mapping

type Mapping struct {
	destinationFieldName string

	sourceField  string
	sourceMethod string
}

func (m Mapping) Destination() string {
	return m.destinationFieldName
}

func (m Mapping) Source() string {
	if m.sourceField != "" {
		return m.sourceField
	}

	return m.sourceMethod
}

func (m Mapping) IsMethod() bool {
	return m.sourceMethod != ""
}
