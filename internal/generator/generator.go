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
	imports := combineImports(srcTypes, dstType)
	mapper := newMapper(srcTypes, dstType)

	g.writer.WritePackage(outputPkg)
	g.writer.WriteImports(imports)
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

func combineImports(srcTypes []*ParsedStruct, dstType *ParsedStruct) Imports {
	imports := dstType.Imports
	imports = append(imports, dstType.Package.ToImport())
	for _, s := range srcTypes {
		imports = append(imports, s.Package.ToImport())
		imports = append(imports, s.Imports...)
	}
	return imports
}
