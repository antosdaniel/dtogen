package struct_with_sub_types_test

import (
	"bytes"
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
)

type StructWithSubTypes struct {
	A string

	// Some system types as registered by default
	Buffer bytes.Buffer
	Time   time.Time

	// Custom types
	RegisteredValueType    _misc.RegisteredValueType
	NonRegisteredValueType NonRegisteredValueType
}
