package generator

type Input struct {
	// PathToSource Import path to package in which DTO is present.
	PathToSource string
	// TypeName Name of DTO in the source.
	TypeName string
	// RenameTypeTo Desired name of a DTO. If empty, original name will be used.
	RenameTypeTo string

	// If IncludeAllParsedFields is set to true, all fields, no matter if they are present in Fields, will be included.
	IncludeAllParsedFields bool
	// Fields Specifies which fields should be included in new DTO, and what they should be renamed to.
	Fields          FieldsInput
	RegisteredTypes RegisteredTypesInput

	// OutputPackage Name of a package that result DTO should belong to.
	OutputPackage string
}

func (i Input) desiredTypeName() string {
	if i.RenameTypeTo != "" {
		return i.RenameTypeTo
	}

	return i.TypeName
}

type FieldsInput []FieldInput

type FieldInput struct {
	Name     string
	RenameTo string
}

func (fs FieldsInput) findByOriginalName(name string) (FieldInput, bool) {
	for _, f := range fs {
		if f.Name == name {
			return f, true
		}
	}

	return FieldInput{}, false
}

type RegisteredTypesInput []RegisteredTypeInput

type RegisteredTypeInput struct {
	ImportPath string
	TypeName   string
}