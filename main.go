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
	src = flag.StringP("src", "s", "", "Path to source object, based on which generation will happen. Value should follow <import path>.<type> pattern. For example \"net/http.Response\". You can provide multiple source object separated by comma. Absolute and relative import paths are supported.")
	dst = flag.StringP("dst", "d", "", "Path to destination package, where generated object will be placed. Value should follow <import path>.<type> pattern. For example \"./net/dto.ResponseDTO\". Absolute and relative import paths are supported.")
	out = flag.StringP("out", "o", "", "Path to file, to which generated code will be saved. If empty, stdout will be used.")

	help = flag.BoolP("help", "h", false, "Show help prompt.")
)

func main() {
	flag.Parse()
	output, err := run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if output != "" {
		fmt.Println(output)
	}
}

func run(args []string) (string, error) {
	if len(args) <= 1 || *help {
		return printHelp(), nil
	}

	input, err := getInput()
	if err != nil {
		return "", fmt.Errorf("invalid argument: %s", err)
	}

	result, err := generator.New(parser.New(), writer.New()).Generate(*input)
	if err != nil {
		return "", fmt.Errorf("generation failed: %s\n", err)
	}

	if *out == "" {
		return result.String(), nil
	}

	return "", saveToFile(result.String())
}

func getInput() (*generator.Input, error) {
	if *src == "" {
		return nil, fmt.Errorf("src: can not be empty")
	}
	srcs := strings.Split(*src, ",")
	srcInput := make([]generator.TypeInput, 0)
	for _, s := range srcs {
		r, err := getTypeInput(s)
		if err != nil {
			return nil, fmt.Errorf("src: %w", err)
		}
		srcInput = append(srcInput, r)
	}

	if *dst == "" {
		return nil, fmt.Errorf("dst: can not be empty")
	}
	dstInput, err := getTypeInput(*dst)
	if err != nil {
		return nil, fmt.Errorf("dst: %w", err)
	}

	return &generator.Input{
		Src: srcInput,
		Dst: dstInput,
	}, nil
}

func getTypeInput(in string) (generator.TypeInput, error) {
	parts := strings.Split(in, ".")
	if len(parts) < 2 {
		return generator.TypeInput{}, fmt.Errorf("requires \"<import path>.<type>\" pattern")
	}

	return generator.TypeInput{
		// Import path can use dots as well. We only treat last part as a type.
		ImportPath: strings.Join(parts[0:len(parts)-1], "."),
		Type:       parts[len(parts)-1],
	}, nil
}

func saveToFile(content string) error {
	err := os.WriteFile(*out, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("could not save file: %s\n", err)
	}
	return nil
}

func printHelp() string {
	return strings.TrimSpace(fmt.Sprintf(`
godtogen --src <value> --dst <value> [--out <value>] [-rename <value>] [fields]

  [fields] should follow a pattern of <name>[=rename][.type]
      <name> determines name of a field on source object. Required.
      [=rename] can be used to rename field. Optional.
      [.type] can be used to override field type. Optional.
      For example: "Foo=Bar.MyType" will rename field "Foo" to "Bar", and change its type to "MyType".
      If no fields are selected, all fields from the source will be used.

%s
`, flag.CommandLine.FlagUsages()))
}
