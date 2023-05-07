package generator

import (
	"go/ast"
	"strings"
)

type ParsedStruct struct {
	Name    string
	Fields  ParsedFields
	Methods Methods
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

type Methods []Method

type Method struct {
	name string
}

func NewMethod(name string) Method {
	return Method{name: name}
}

func (ms Methods) findExportedMethod(names ...string) (Method, bool) {
	for _, name := range names {
		for _, method := range ms {
			if !method.IsExported() {
				continue
			}
			if strings.ToLower(name) != strings.ToLower(method.name) {
				continue
			}

			return method, true
		}
	}

	return Method{}, false
}

func (m Method) IsExported() bool {
	return ast.IsExported(m.name)
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

func (p Package) Name() string {
	return p.name
}

type ParsedFunctions []ParsedFunction

type ParsedFunction struct {
	name string
	body string
}

func NewParsedFunction(name, body string) ParsedFunction {
	return ParsedFunction{name: name, body: body}
}

func (pf ParsedFunction) Code() string {
	return pf.body
}
