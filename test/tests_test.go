package test_test

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"github.com/antosdaniel/dtogen/internal/parser"
	"github.com/antosdaniel/dtogen/internal/writer"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		testdata string
		input    generator.Input
	}{
		{
			testdata: "rename_struct",
			input: generator.Input{
				RenameTypeTo:           "RenamedDTO",
				IncludeAllParsedFields: true,
			},
		},
		{
			testdata: "struct_with_base_types",
			input: generator.Input{
				IncludeAllParsedFields: true,
			},
		},
		{
			testdata: "struct_with_override_type",
			input: generator.Input{
				IncludeAllParsedFields: true,
				Fields: generator.FieldsInput{
					{
						Name:           "Policy",
						OverrideTypeTo: "PolicyDTO",
					},
				},
			},
		},
		{
			testdata: "struct_with_renamed_fields",
			input: generator.Input{
				IncludeAllParsedFields: true,
				Fields: []generator.FieldInput{
					{Name: "B", RenameTo: "X"},
					{Name: "C", RenameTo: "Y"},
				},
			},
		},
		{
			testdata: "struct_with_sub_types",
			input: generator.Input{
				IncludeAllParsedFields: true,
			},
		},
		{
			testdata: "struct_with_unexported_field",
			input: generator.Input{
				IncludeAllParsedFields: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testdata, func(t *testing.T) {
			tc.input.TypeName = "DTO"
			tc.input.FilePath = fmt.Sprintf("./testdata/%s/input.go", tc.testdata)
			tc.input.OutputPackage = fmt.Sprintf("%s_test", tc.testdata)
			g := generator.New(parser.New(), writer.New())

			result, err := g.Generate(tc.input)

			if assert.NoError(t, err) {
				e := expected(t, fmt.Sprintf("./testdata/%s/output_test.go", tc.testdata))
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
