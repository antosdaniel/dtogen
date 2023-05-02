package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antosdaniel/dtogen/internal/generator"
	"github.com/antosdaniel/dtogen/internal/parser"
	"github.com/antosdaniel/dtogen/internal/writer"
	flag "github.com/spf13/pflag"
)

var (
	src = flag.StringP("src", "s", "", "Path to source object, based on which generation will happen.\nValue should follow <import path>.<type> pattern. For example \"net/http.Response\". Absolute and relative import paths are supported.")
	dst = flag.StringP("dst", "d", "", "Path to destination package, where generated object will be placed.\nAbsolute and relative import paths are supported.")
	out = flag.StringP("out", "o", "", "Path to file, to which generated code will be saved. If empty, stdout will be used.")

	rename    = flag.StringP("rename", "r", "", "Desired name of generated object. If empty, source name will be used.")
	allFields = flag.Bool("all-fields", false, "All fields of source object will be used. You can still use it with [fields] arguments, if you want to rename field, or change its type.")

	help = flag.BoolP("help", "h", false, "Show help prompt.")
)

func main() {
	flag.Parse()
	if len(os.Args) <= 1 || *help {
		printHelp()
		os.Exit(0)
	}

	input, err := getInput()
	if err != nil {
		fmt.Printf("Invalid argument: %s\n", err)
		os.Exit(1)
	}

	result, err := generator.New(parser.New(), writer.New()).Generate(*input)
	if err != nil {
		fmt.Printf("Generation failed: %s\n", err)
		os.Exit(1)
	}

	if *out == "" {
		fmt.Println(result.String())
		os.Exit(0)
	}

	saveToFile(result.String())
}

func getInput() (*generator.Input, error) {
	if *src == "" {
		return nil, fmt.Errorf("src can not be empty")
	}
	srcParts := strings.Split(*src, ".")
	if len(srcParts) < 2 {
		return nil, fmt.Errorf("src should follow <import path>.<type> pattern")
	}
	srcPkgPath := strings.Join(srcParts[0:len(srcParts)-1], ".")
	srcType := srcParts[len(srcParts)-1]

	if *dst == "" {
		return nil, fmt.Errorf("dst can not be empty")
	}

	fields := getFields(flag.Args())

	return &generator.Input{
		PackagePath:            srcPkgPath,
		TypeName:               srcType,
		RenameTypeTo:           *rename,
		IncludeAllParsedFields: *allFields || len(fields) == 0,
		Fields:                 fields,
		OutputPackagePath:      *dst,
	}, nil
}

func getFields(args []string) generator.FieldsInput {
	const renameChar = '='
	const typeChar = '.'

	fields := generator.FieldsInput{}
	for _, a := range args {
		if a == "" {
			continue
		}

		hasRename := strings.ContainsRune(a, renameChar)
		hasType := strings.ContainsRune(a, typeChar)
		if !hasRename && !hasType {
			fields = append(fields, generator.FieldInput{Name: a})
			continue
		}

		parts := strings.FieldsFunc(a, func(r rune) bool {
			return r == renameChar || r == typeChar
		})
		if hasRename && !hasType {
			fields = append(fields, generator.FieldInput{Name: parts[0], RenameTo: parts[1]})
			continue
		}
		if !hasRename && hasType {
			fields = append(fields, generator.FieldInput{Name: parts[0], OverrideTypeTo: parts[1]})
			continue
		}
		fields = append(fields, generator.FieldInput{Name: parts[0], RenameTo: parts[1], OverrideTypeTo: parts[2]})
	}
	return fields
}

func saveToFile(content string) {
	err := os.WriteFile(*out, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Could not save file: %s\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("godtogen --src <value> --dst <value> [--out <value>] [-rename <value>] [fields]")
	fmt.Println()
	fmt.Println(`  [fields] should follow a pattern of <name>[=rename][.type]
      <name> determines name of a field on source object. Required.
      [=rename] can be used to rename field. Optional.
      [.type] can be used to override field type. Optional.
      For example: "Foo=Bar.MyType" will rename field "Foo" to "Bar", and change its type to "MyType".
      If no fields are selected, all fields from the source will be used.`)
	fmt.Println()
	flag.PrintDefaults()
}
