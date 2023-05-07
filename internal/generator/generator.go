package generator

import "strings"

type Generator struct {
	parser Parser
	writer Writer
}

func New(parser Parser, writer Writer) *Generator {
	return &Generator{
		parser: parser,
		writer: writer,
	}
}

func (g *Generator) Generate(input Input) (Generated, error) {
	var srcTypes []*ParsedStruct
	for _, src := range input.Src {
		t, err := g.parser.GetStruct(src.ImportPath, src.Type)
		if err != nil {
			return nil, err
		}
		srcTypes = append(srcTypes, t)
	}

	dstType, err := g.parser.GetStruct(input.Dst.ImportPath, input.Dst.Type)
	if err != nil {
		return nil, err
	}

	outputPkg := getOutputPkg(input.OutputPkgPath)
	mapper := newMapper(srcTypes, dstType)

	g.writer.WritePackage(outputPkg)
	g.writer.WriteEmptyLine()
	g.writer.WriteMapper(mapper, outputPkg)

	return g.writer, nil
}

func getOutputPkg(outputPkgPath string) string {
	parts := strings.Split(outputPkgPath, "/")
	if len(parts) <= 1 {
		return outputPkgPath
	}

	return parts[len(parts)-1]
}
