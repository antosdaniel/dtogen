package test_test

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		testdata string
		input    internal.Input
	}{
		{
			testdata: "struct_with_base_types",
			input: internal.Input{
				TypeName:               "StructWithBaseTypes",
				IncludeAllParsedFields: true,
			},
		},
		{
			testdata: "rename_struct",
			input: internal.Input{
				TypeName:               "RenameStruct",
				RenameTypeTo:           "RenameStructToSomethingElse",
				IncludeAllParsedFields: true,
			},
		},
		{
			testdata: "struct_with_renamed_fields",
			input: internal.Input{
				TypeName:               "StructWithRenamedFields",
				IncludeAllParsedFields: true,
				Fields: []internal.FieldInput{
					{Name: "B", RenameTo: "X"},
					{Name: "C", RenameTo: "Y"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testdata, func(t *testing.T) {
			testCase.input.PathToSource = fmt.Sprintf("./testdata/%s/input.go", testCase.testdata)
			testCase.input.OutputPackage = fmt.Sprintf("%s_test", testCase.testdata)

			result, err := internal.Generate(testCase.input)

			if assert.NoError(t, err) {
				e := expected(t, fmt.Sprintf("./testdata/%s/output_test.go", testCase.testdata))
				assert.Equal(t, e, result.String())
			}
		})
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
