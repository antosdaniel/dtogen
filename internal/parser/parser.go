package parser

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type Parser struct {
	pkg *packages.Package
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) LoadPackage(importPath string) (generator.Parser, error) {
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

func (p *Parser) GetStruct(typeName string) (*generator.ParsedStruct, error) {
	if p.pkg == nil {
		return nil, fmt.Errorf("no package loaded")
	}

	for _, file := range p.pkg.Syntax {
		for _, d := range file.Decls {
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
					return parseStruct(t, typeName, file)
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

func parseStruct(structType *ast.StructType, typeName string, file *ast.File) (*generator.ParsedStruct, error) {
	imports := make(generator.Imports, len(file.Imports))
	for _, i := range file.Imports {
		imports = append(imports, generator.NewImport(i))
	}

	result := &generator.ParsedStruct{
		Name:    typeName,
		Imports: imports,
		Fields:  generator.ParsedFields{},
	}

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			result.Fields = append(result.Fields, generator.NewParsedField(name.Name, field.Type))
		}
	}

	return result, nil
}
