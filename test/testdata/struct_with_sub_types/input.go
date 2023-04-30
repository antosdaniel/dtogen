package struct_with_sub_types

import (
	"bytes"
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
	mytime "time"
)

type DTO struct {
	// Some system types as registered by default
	Buffer bytes.Buffer
	Time   *time.Time
	MyTime mytime.Time

	// Custom types
	RegisteredValueType    _misc.RegisteredType
	NonRegisteredValueType _misc.NonRegisteredType
}
