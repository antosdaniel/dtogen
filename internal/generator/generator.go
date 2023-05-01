package generator

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
	p, err := g.parser.LoadPackage(input.PackagePath)
	if err != nil {
		return nil, err
	}
	parsed, err := p.GetStruct(input.TypeName)
	if err != nil {
		return nil, err
	}

	// TODO: error on not found names
	fields := prepareFields(input.IncludeAllParsedFields, input.Fields, parsed.Fields, parsed.Imports)
	imports := fields.RequiredImports()
	// If we are generating mapper, then we will need to import source package.
	if input.generateMapper() {
		imports = append(imports, parsed.Package.ToImport())
	}

	g.writer.WritePackage(input.OutputPackage)
	g.writer.WriteEmptyLine()
	if len(imports) > 0 {
		g.writer.WriteImports(imports)
		g.writer.WriteEmptyLine()
	}
	g.writer.WriteStruct(input.desiredTypeName(), fields)
	if input.generateMapper() {
		g.writer.WriteEmptyLine()
		g.writer.WriteMapper(prepareMapper(*parsed, fields, input.desiredTypeName()))
	}

	return g.writer, nil
}
