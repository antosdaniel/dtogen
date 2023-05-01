package struct_mapper_with_handwritten_mappings

import (
	"time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
)

type Input struct {
	ID     string
	Policy _misc.Policy

	CreatedAt time.Time
	DeletedAt *time.Time
}
