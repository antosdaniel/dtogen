package generator

type Parser interface {
	GetStruct(importPath, typeName string) (*ParsedStruct, error)
	GetFunctions(importPath string) (ParsedFunctions, Imports, error)
}

type Writer interface {
	String() string
	WritePackage(pkg string)
	WriteImports(imports Imports)
	WriteMapper(mapper Mapper, outputPkg string)
}

type Generated interface {
	String() string
}
