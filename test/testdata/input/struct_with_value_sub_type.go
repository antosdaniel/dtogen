package input

import (
	"github.com/antosdaniel/dtogen/test/testdata/value_type"
)

type StructWithValueSubType struct {
	A     string
	Value value_type.ValueType
}
