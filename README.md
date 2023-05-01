# dtogen [![Tests](https://github.com/antosdaniel/dtogen/actions/workflows/test.yml/badge.svg)](https://github.com/antosdaniel/dtogen/actions) [![Security](https://github.com/antosdaniel/dtogen/actions/workflows/security.yml/badge.svg)](https://github.com/antosdaniel/dtogen/actions)

Generate DTOs with mappers based on existing Go types.

## Install

```sh
go install github.com/antosdaniel/dtogen/cmd/godtogen
```

## Example

### Create DTO and mapper for built-in type

`Body/ResponseBody` renames `Body` field to `ResponseBody`

```sh
godtogen -pkg-path net/http -type-name Response Status StatusCode Body/ResponseBody Request
```

```go
package http

import (
	"io"
	"net/http"
)

type Response struct {
	Status       string
	StatusCode   int
	ResponseBody io.ReadCloser
}

func NewResponse(src http.Response) Response {
	return Response{
		Status:       src.Status,
		StatusCode:   src.StatusCode,
		ResponseBody: src.Body,
	}
}
```

### Create DTO and mapper using auto-detected getters

```sh
godtogen -pkg-path ./test/testdata/struct_mapper_with_getters -type-name Input
```

```go
package struct_mapper_with_getters

import (
    "time"

    "github.com/antosdaniel/dtogen/test/testdata/_misc"
    "github.com/antosdaniel/dtogen/test/testdata/struct_mapper_with_getters"
)

type Input struct {
    Id        string
    Metadata  _misc.CustomType
    CreatedAt time.Time
    DeletedAt *time.Time
}

func NewInput(src struct_mapper_with_getters.Input) Input {
    return Input{
        Id:        src.ID(),
        Metadata:  src.Metadata,
        CreatedAt: src.GetCreatedAt(),
        DeletedAt: src.DeletedAt(),
    }
}
```

### Change field type, and use handwritten mapping functions

You might want to change certain fields' type. Additionaly, you might want to control how these types are created. Any `New<Field>` (exported or unexported) in output package will be preserved and used.

```sh
godtogen -pkg-path ./test/testdata/struct_mapper_with_handwritten_mappings -out-pkg ./test/testdata/struct_mapper_with_handwritten_mappings/output -type-name Input -all-fields Policy//SimplePolicy
```

`Policy//SimplePolicy` instructs that `Policy` field should be of `SimplePolicy` type.

```go
package output

import (
	"time"

	"github.com/antosdaniel/dtogen/test/testdata/_misc"
	"github.com/antosdaniel/dtogen/test/testdata/struct_mapper_with_handwritten_mappings"
)

type Input struct {
	ID        string
	Policy    SimplePolicy
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewInput(src struct_mapper_with_handwritten_mappings.Input) Input {
	return Input{
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
	}
}
```
