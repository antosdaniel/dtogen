package struct_mapper_with_getters_test

import (
	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_mapper_with_getters"
	"time"
)

type DTO struct {
	Id        string
	Metadata  _misc.CustomType
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewDTO(src struct_mapper_with_getters.Input) DTO {
	return DTO{
		Id:        src.ID(),
		Metadata:  src.Metadata,
		CreatedAt: src.GetCreatedAt(),
		DeletedAt: src.DeletedAt(),
	}
}
