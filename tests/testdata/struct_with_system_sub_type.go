package testdata

import (
	"bytes"
	"time"
)

type StructWithSystemSubType struct {
	A        string
	SubTypeA bytes.Buffer
	SubTypeB time.Time
}
