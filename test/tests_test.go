package test_test

import (
	"os"
	"testing"

	"github.com/antosdaniel/mappergen/internal/generator"
	"github.com/antosdaniel/mappergen/internal/parser"
	"github.com/antosdaniel/mappergen/internal/writer"
	"github.com/stretchr/testify/assert"
)

func TestScenarios(t *testing.T) {
	testCases := []struct {
		name     string
		input    generator.Input
		expected string
	}{
		{
			name: "simple_struct",
			input: generator.Input{
				Src: []generator.TypeInput{
					{
						ImportPath: "./testdata/simple_struct",
						Type:       "Src",
					},
				},
				Dst: generator.TypeInput{
					ImportPath: "./testdata/simple_struct",
					Type:       "Dst",
				},
				OutputPkgPath: "./testdata/simple_struct",
			},
			expected: "./testdata/simple_struct/mapper.go",
		},
		{
			name: "nested_types",
			input: generator.Input{
				Src: []generator.TypeInput{
					{
						ImportPath: "./testdata/nested_types/src",
						Type:       "Payslip",
					},
				},
				Dst: generator.TypeInput{
					ImportPath: "./testdata/nested_types/dst",
					Type:       "Payslip",
				},
				OutputPkgPath: "./testdata/nested_types",
			},
			expected: "./testdata/nested_types/mapper.go",
		},
		{
			name: "import_paths",
			input: generator.Input{
				Src: []generator.TypeInput{
					{
						ImportPath: "./testdata/import_paths/src",
						Type:       "Src",
					},
				},
				Dst: generator.TypeInput{
					ImportPath: "./testdata/import_paths/dst",
					Type:       "Dst",
				},
				OutputPkgPath: "./testdata/import_paths",
			},
			expected: "./testdata/import_paths/mapper.go",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := generator.New(parser.New(), writer.New())

			result, err := g.Generate(tc.input)

			if assert.NoError(t, err) {
				assert.Equal(t, getExpected(t, tc.expected), result.String())
			}
		})
	}
}

func getExpected(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to load expected result from %q", path)
		return ""
	}
	return string(content)
}
