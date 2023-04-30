package struct_with_sub_types

import (
	"bytes"
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
	mytime "time"
)

type Input struct {
	// Some system types as registered by default
	Buffer bytes.Buffer
	Time   *time.Time
	MyTime mytime.Time

	// Custom types
	CustomType _misc.CustomType
}
