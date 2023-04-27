package test_test

import (
	"github.com/antosdaniel/dtogen/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		input    internal.Input
		expected string
	}{
		{
			input: internal.Input{
				PathToSource:           "./testdata/input/struct_with_base_types.go",
				TypeName:               "StructWithBaseTypes",
				IncludeAllParsedFields: true,
				OutputPackage:          "output",
			},
			expected: "./testdata/output/struct_with_base_types.go",
		},
		{
			input: internal.Input{
				PathToSource:           "./testdata/input/rename_struct.go",
				TypeName:               "RenameStruct",
				RenameTypeTo:           "RenameStructToSomethingElse",
				IncludeAllParsedFields: true,
				OutputPackage:          "output",
			},
			expected: "./testdata/output/rename_struct.go",
		},
		{
			input: internal.Input{
				PathToSource:           "./testdata/input/struct_with_renamed_fields.go",
				TypeName:               "StructWithRenamedFields",
				IncludeAllParsedFields: true,
				Fields: []internal.FieldInput{
					{Name: "B", RenameTo: "X"},
					{Name: "C", RenameTo: "Y"},
				},
				OutputPackage: "output",
			},
			expected: "./testdata/output/struct_with_renamed_fields.go",
		},
	}

	for _, testCase := range testCases {
		result, err := internal.Generate(testCase.input)

		if assert.NoError(t, err) {
			assert.Equal(t, expected(t, testCase.expected), result.String())
		}
	}
}

func expected(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to load expected result from %q", path)
		return ""
	}
	return string(content)
}
