package struct_with_sub_types_test

import (
	"bytes"
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
	mytime "time"
)

type DTO struct {
	Buffer     bytes.Buffer
	Time       *time.Time
	MyTime     mytime.Time
	CustomType _misc.CustomType
}
