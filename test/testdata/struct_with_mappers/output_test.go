package struct_with_mappers_test

import (
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_with_mappers"
	"time"
)

type DTO struct {
	ID        string
	Metadata  _misc.CustomType
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewDTO(src struct_with_mappers.Input) DTO {
	return DTO{
		ID:        src.ID,
		Metadata:  src.Metadata,
		CreatedAt: src.CreatedAt,
		DeletedAt: src.DeletedAt,
	}
}
