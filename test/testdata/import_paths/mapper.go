package import_paths

import (
	"github.com/antosdaniel/dtogen/test/testdata/import_paths/dst"
	"github.com/antosdaniel/dtogen/test/testdata/import_paths/src"
)

func ToDst(src src.Src) dst.Dst {
	return dst.Dst{
		Bool:   src.Bool,
		Int:    src.Int,
		String: src.String,
		Time:   src.Time,
	}
}
