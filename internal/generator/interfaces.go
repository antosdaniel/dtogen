package generator

type Parser interface {
	LoadPackage(importPath string) (Parser, error)
	GetStruct(typeName string) (*ParsedStruct, error)
}

type Writer interface {
	String() string
	In()
	Out()
	Write(code string)
	WriteLine(code string)
	WriteEmptyLine()
	WritePackage(pkg string)
	WriteImports(imports Imports)
	WriteStruct(structName string, fields Fields)
}

type Generated interface {
	String() string
}
