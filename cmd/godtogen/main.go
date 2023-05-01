package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/antosdaniel/dtogen/internal/generator"
	"github.com/antosdaniel/dtogen/internal/parser"
	"github.com/antosdaniel/dtogen/internal/writer"
)

var (
	packagePath  = flag.String("pkg-path", "", "Import path to package in which source object is present.")
	typeName     = flag.String("type-name", "", "Name of DTO in the source.")
	renameTypeTo = flag.String("rename-type-to", "", "Desired name of a DTO. If empty, original name will be used.")

	includeAllParsedFields = flag.Bool("all-fields", false, "If set to true, all fields, no matter if they are present in Fields, will be included.")
	skipMapper             = flag.Bool("skip-mapper", false, "If set to true, mapper will not be generated.")

	outputPackagePath = flag.String("out-pkg", "", "Import path to where DTO will be generated.")
	outputFile        = flag.String("o", "", "File to which generated Go file will be saved. If empty, stdout will be used.")
)

func main() {
	flag.Parse()
	if len(os.Args) <= 1 {
		printHelp()
		os.Exit(0)
	}

	fields := generator.FieldsInput{}
	for _, a := range flag.Args() {
		parts := strings.Split(a, "/")
		var name, rename, _type string
		if len(parts) >= 1 {
			name = parts[0]
		}
		if len(parts) >= 2 {
			rename = parts[1]
		}
		if len(parts) >= 3 {
			_type = parts[2]
		}
		fields = append(fields, generator.FieldInput{Name: name, RenameTo: rename, OverrideTypeTo: _type})
	}

	input := generator.Input{
		PackagePath:            *packagePath,
		TypeName:               *typeName,
		RenameTypeTo:           *renameTypeTo,
		IncludeAllParsedFields: *includeAllParsedFields || len(fields) == 0,
		Fields:                 fields,
		SkipMapper:             *skipMapper,
		OutputPackagePath:      *outputPackagePath,
	}
	result, err := generator.New(parser.New(), writer.New()).Generate(input)
	if err != nil {
		fmt.Printf("Error during generating: %s\n", err)
		os.Exit(1)
	}

	if *outputFile == "" {
		fmt.Println(result.String())
		os.Exit(0)
	}

	err = os.WriteFile(*outputFile, []byte(result.String()), 0644)
	if err != nil {
		fmt.Printf("Could not save file: %s\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`Generate DTOs with mappers based on existing Go types.

godtogen [flags] [fields]

Generate with all fields:
godtogen -pkg-path github.com/antosdaniel/dtogen/test/testdata/struct_with_base_types -type-name Input -out-pkg github.com/antosdaniel/dtogen

Generate with specific fields only:
godtogen -pkg-path github.com/antosdaniel/dtogen/test/testdata/struct_with_base_types -type-name Input -out-pkg github.com/antosdaniel/dtogen Bool Int Uint/Foo

`)
	flag.Usage()
}
