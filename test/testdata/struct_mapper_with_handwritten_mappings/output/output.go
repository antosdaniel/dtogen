package output

import (
	"strings"
	"time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_mapper_with_handwritten_mappings"
)

type DTO struct {
	ID        string
	Policy    SimplePolicy
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewDTO(src struct_mapper_with_handwritten_mappings.Input) DTO {
	return DTO{
		ID:        src.ID,
		Policy:    NewPolicy(src.Policy),
		CreatedAt: src.CreatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func NewPolicy(src _misc.Policy) SimplePolicy {
	return SimplePolicy{
		ID:   src.ID,
		Name: src.Name,
		// strings.Join is used to check if imports will be preserved
		AuthorFullName: strings.Join([]string{src.AuthorFirstName, src.AuthorLastName}, " - "),
	}
}
