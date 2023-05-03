package struct_with_sub_types

import (
	"bytes"
	"go/ast"
	"time"
	mytime "time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_with_sub_types/key"
	"github.com/antosdaniel/dtogen/test/testdata/struct_with_sub_types/value"
)

type Input struct {
	// Some system types
	Buffer bytes.Buffer
	Time   *time.Time
	MyTime mytime.Time
	Types  []ast.Expr

	// Custom types
	CustomType _misc.CustomType
	CustomMap  map[key.Key]value.Value
}
