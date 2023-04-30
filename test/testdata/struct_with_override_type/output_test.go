package struct_with_override_type_test

import (
	"time"
)

type DTO struct {
	ID        string
	Name      string
	Policy    PolicyDTO
	CreatedAt *time.Time
}
