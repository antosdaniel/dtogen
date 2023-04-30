package generator

import "go/ast"

type ParsedStruct struct {
	Name    string
	Fields  ParsedFields
	Imports Imports
	Package Package
}

type ParsedFields []ParsedField

type ParsedField struct {
	name  string
	_type ast.Expr
}

func NewParsedField(name string, _type ast.Expr) ParsedField {
	return ParsedField{name: name, _type: _type}
}

func (f ParsedField) Name() string {
	return f.name
}

func (f ParsedField) Type() ast.Expr {
	return f._type
}

type Package struct {
	path string
	name string
}

func NewPackage(path, name string) Package {
	return Package{
		path: path,
		name: name,
	}
}

func (p Package) ToImport() Import {
	return Import{
		path: p.path,
		name: p.name,
	}
}
