package test_test

import (
	"github.com/antosdaniel/dtogen/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScenarios(t *testing.T) {
	t.Run("struct with base types", func(t *testing.T) {
		result, err := internal.Generate(internal.Input{
			PathToSource:           "./testdata/input/struct_with_base_types.go",
			TypeName:               "StructWithBaseTypes",
			IncludeAllParsedFields: true,
			OutputPackage:          "output",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, expected(t, "./testdata/output/struct_with_base_types.go"), result.String())
		}
	})
	t.Run("rename struct", func(t *testing.T) {
		result, err := internal.Generate(internal.Input{
			PathToSource:           "./testdata/input/rename_struct.go",
			TypeName:               "RenameStruct",
			RenameTypeTo:           "RenameStructToSomethingElse",
			IncludeAllParsedFields: true,
			OutputPackage:          "output",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, expected(t, "./testdata/output/rename_struct.go"), result.String())
		}
	})
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
