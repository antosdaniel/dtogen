package struct_with_override_type

import (
	"time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
)

type Input struct {
	ID        string
	Name      string
	Policy    _misc.CustomType
	CreatedAt *time.Time
}
