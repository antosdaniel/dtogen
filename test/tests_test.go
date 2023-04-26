package test_test

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal"
	"testing"
)

func TestScenarios(t *testing.T) {
	t.Run("struct with base types", func(t *testing.T) {
		result, err := internal.Parse("./testdata/input/struct_with_base_types.go", "StructWithBaseTypes")

		fmt.Println(result, err)
	})
}
