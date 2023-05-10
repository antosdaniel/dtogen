# dtogen [![Tests](https://github.com/antosdaniel/mappergen/actions/workflows/test.yml/badge.svg)](https://github.com/antosdaniel/mappergen/actions) [![Coverage](https://coveralls.io/repos/github/antosdaniel/dtogen/badge.svg)](https://coveralls.io/github/antosdaniel/dtogen) [![Go Report Card](https://goreportcard.com/badge/github.com/antosdaniel/mappergen)](https://goreportcard.com/report/github.com/antosdaniel/mappergen) [![Security](https://github.com/antosdaniel/mappergen/actions/workflows/security.yml/badge.svg)](https://github.com/antosdaniel/mappergen/actions)

Generate DTOs with mappers based on existing Go types.

## Install

```sh
go install github.com/antosdaniel/mappergen/cmd/godtogen
```

## Usage

```
godtogen --src <value> --dst <value> [--out <value>] [-rename <value>] [fields]

  [fields] should follow a pattern of <name>[=rename][.type]
      <name> determines name of a field on source object. Required.
      [=rename] can be used to rename field. Optional.
      [.type] can be used to override field type. Optional.
      For example: "Foo=Bar.MyType" will rename field "Foo" to "Bar", and change its type to "MyType".
      If no fields are selected, all fields from the source will be used.

      --all-fields      All fields of source object will be used. You can still use it with [fields] arguments, if you want to rename field, or change its type.
  -d, --dst string      Path to destination package, where generated object will be placed.
                        Absolute and relative import paths are supported.
  -h, --help            Show help prompt.
  -o, --out string      Path to file, to which generated code will be saved. If empty, stdout will be used.
  -r, --rename string   Desired name of generated object. If empty, source name will be used.
  -s, --src string      Path to source object, based on which generation will happen.
                        Value should follow <import path>.<type> pattern. For example "net/http.Response". Absolute and relative import paths are supported.
```

## Examples

### Create DTO and mapper for built-in type, with selected fields

Create DTO for type `Response` from package `net/http`. Destination package is set to relative `./my/dto`. We are picking 4 specific fields, and specifying that `Body` field should be renamed to `ResponseBody`.

```sh
godtogen --src net/http.Response --dst ./my/dto Status StatusCode Body=ResponseBody Request
```

```go
package dto

import (
	"io"
	"net/http"
)

type Response struct {
	Status       string
	StatusCode   int
	ResponseBody io.ReadCloser
	Request      *Request
}

func NewResponse(src http.Response) Response {
	return Response{
		Status:       src.Status,
		StatusCode:   src.StatusCode,
		ResponseBody: src.Body,
		Request:      src.Request,
	}
}
```

### Map fields using getters

When source fields is not exported, mapper will automatically look for getters.

```sh
godtogen --src ./test/testdata/struct_mapper_with_getters.Input --dst dto
```

```go
package dto

import (
	"time"

	"github.com/antosdaniel/mappergen/test/testdata/_misc"
	"github.com/antosdaniel/mappergen/test/testdata/struct_mapper_with_getters"
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
