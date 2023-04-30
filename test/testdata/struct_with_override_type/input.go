package struct_with_override_type

import (
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
)

type DTO struct {
	ID        string
	Name      string
	Policy    _misc.CustomType
	CreatedAt *time.Time
}
