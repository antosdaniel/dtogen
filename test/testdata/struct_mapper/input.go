package struct_mapper

import (
	"time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
)

type Input struct {
	ID       string
	Metadata _misc.CustomType

	CreatedAt time.Time
	DeletedAt *time.Time
}
