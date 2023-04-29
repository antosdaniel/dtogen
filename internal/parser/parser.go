package parser

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type Parser struct {
	pkg *packages.Package
}

func New(importPath string) (*Parser, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedImports,
	}, importPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package not found: %q", importPath)
	}

	return &Parser{pkg: pkgs[0]}, nil
}

func (p *Parser) GetStruct(typeName string) (*Struct, error) {
	for _, s := range p.pkg.Syntax {
		for _, d := range s.Decls {
			genDecl, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if typeName != typeSpec.Name.Name {
					continue
				}

				switch typeSpec.Type.(type) {
				case *ast.StructType:
					t := typeSpec.Type.(*ast.StructType)
					return parseStruct(t, typeName)
				default:
					return nil, fmt.Errorf(
						"specified declaration is not a struct: %T", typeSpec.Type,
					)
				}
			}
		}
	}

	return nil, fmt.Errorf("type not found: %q", typeName)
}

func parseStruct(structType *ast.StructType, typeName string) (*Struct, error) {
	result := &Struct{
		Name:   typeName,
		Fields: Fields{},
	}

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			result.Fields = append(result.Fields, Field{
				Name: name.Name,
				Type: field.Type,
			})
		}
	}

	return result, nil
}

type Struct struct {
	Name   string
	Fields Fields
}

type Fields []Field

type Field struct {
	Name string
	Type ast.Expr
}
