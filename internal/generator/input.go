package generator

// Input
//
// Required:
// - PackagePath
// - TypeName
// - IncludeAllParsedFields or Fields. Both can be present, if you want to include all fields, and rename some of them.
// - OutputPackage
type Input struct {
	// PackagePath Import path to package in which DTO is present.
	PackagePath string
	// TypeName Name of DTO in the source.
	TypeName string
	// RenameTypeTo Desired name of a DTO. If empty, original name will be used.
	RenameTypeTo string

	// If IncludeAllParsedFields is set to true, all fields, no matter if they are present in Fields, will be included.
	IncludeAllParsedFields bool
	// Fields Specifies which fields should be included in new DTO. Each field can be renamed, or have override type.
	Fields FieldsInput

	// If SkipMapper is set to true, mapper will not be generated.
	SkipMapper bool

	// OutputPackage Name of a package that result DTO should belong to.
	OutputPackage     string
	OutputPackagePath string
}

func (i Input) desiredTypeName() string {
	if i.RenameTypeTo != "" {
		return i.RenameTypeTo
	}

	return i.TypeName
}

func (i Input) generateMapper() bool {
	return !i.SkipMapper
}

type FieldsInput []FieldInput

type FieldInput struct {
	Name           string
	RenameTo       string
	OverrideTypeTo string
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
