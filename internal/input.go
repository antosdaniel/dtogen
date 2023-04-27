package internal

type Input struct {
	// PathToSource Path to file in which DTO is present.
	PathToSource string
	// TypeName Name of DTO in the source.
	TypeName string
	// RenameTypeTo Desired name of a DTO. If empty, original name will be used.
	RenameTypeTo string

	// If IncludeAllParsedFields is set to true, all fields, no matter if they are present in Fields, will be included.
	IncludeAllParsedFields bool
	// Fields Specifies which fields should be included in new DTO, and what they should be renamed to.
	Fields FieldsInput

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
	Name     FieldName
	RenameTo FieldName
}

func (fs FieldsInput) findByOriginalName(name FieldName) (FieldInput, bool) {
	for _, f := range fs {
		if f.Name == name {
			return f, true
		}
	}

	return FieldInput{}, false
}
