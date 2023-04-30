package struct_with_mappers

import (
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
)

type Input struct {
	ID       string
	Metadata _misc.CustomType

	CreatedAt time.Time
	DeletedAt *time.Time
}
