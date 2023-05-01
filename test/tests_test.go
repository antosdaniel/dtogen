package test_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/antosdaniel/dtogen/internal/generator"
	"github.com/antosdaniel/dtogen/internal/parser"
	"github.com/antosdaniel/dtogen/internal/writer"
	"github.com/stretchr/testify/assert"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		testdata string
		input    generator.Input
		expected string
	}{
		{
			testdata: "struct_mapper",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             false,
			},
		},
		{
			testdata: "struct_mapper_with_getters",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             false,
			},
		},
		{
			testdata: "struct_mapper_with_handwritten_mappings",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             false,
				Fields: generator.FieldsInput{
					{
						Name:           "Policy",
						OverrideTypeTo: "SimplePolicy",
					},
				},
				OutputPackage:     "output",
				OutputPackagePath: "./testdata/struct_mapper_with_handwritten_mappings/output",
			},
			expected: "./testdata/struct_mapper_with_handwritten_mappings/output/output.go",
		},
		{
			testdata: "struct_with_base_types",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             true,
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
				SkipMapper: true,
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
				SkipMapper: true,
			},
		},
		{
			testdata: "struct_with_sub_types",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             true,
			},
		},
		{
			testdata: "struct_with_unexported_field",
			input: generator.Input{
				IncludeAllParsedFields: true,
				SkipMapper:             true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testdata, func(t *testing.T) {
			tc.input.TypeName = "Input"
			tc.input.RenameTypeTo = "DTO"
			tc.input.PackagePath = fmt.Sprintf("./testdata/%s", tc.testdata)
			if tc.input.OutputPackage == "" {
				tc.input.OutputPackage = fmt.Sprintf("%s_test", tc.testdata)
			}
			g := generator.New(parser.New(), writer.New())

			result, err := g.Generate(tc.input)

			if assert.NoError(t, err) {
				e := tc.expected
				if e == "" {
					e = fmt.Sprintf("./testdata/%s/output_test.go", tc.testdata)
				}
				assert.Equal(t, expected(t, e), result.String())
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
