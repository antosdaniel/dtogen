package testdata

import (
	"github.com/antosdaniel/dtogen/tests/testdata/value_type"
)

type StructWithValueSubType struct {
	A     string
	Value value_type.ValueType
}
