package struct_with_sub_types_test

import (
	"bytes"
	"go/ast"
	"time"
	mytime "time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_with_sub_types/key"
	"github.com/antosdaniel/dtogen/test/testdata/struct_with_sub_types/value"
)

type DTO struct {
	Buffer     bytes.Buffer
	Time       *time.Time
	MyTime     mytime.Time
	Types      []ast.Expr
	CustomType _misc.CustomType
	CustomMap  map[key.Key]value.Value
}
