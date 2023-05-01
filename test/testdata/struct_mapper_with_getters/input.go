package struct_mapper_with_getters

import (
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"time"
)

type Input struct {
	id       string
	Metadata _misc.CustomType

	createdAt time.Time
	deletedAt *time.Time
}

func (i Input) ID() string {
	return i.id
}

func (i Input) GetCreatedAt() time.Time {
	return i.createdAt
}

func (i Input) DeletedAt() *time.Time {
	return i.deletedAt
}
