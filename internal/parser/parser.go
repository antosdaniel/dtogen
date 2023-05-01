package parser

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type parser struct {
	pkg *packages.Package
}

func New() generator.Parser {
	return &parser{}
}

func (p *parser) LoadPackage(importPath string) (generator.Parser, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedImports,
	}, importPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package not found: %q", importPath)
	}

	return &parser{pkg: pkgs[0]}, nil
}

func (p *parser) GetStruct(structName string) (*generator.ParsedStruct, error) {
	if p.pkg == nil {
		return nil, fmt.Errorf("load package first")
	}

	parsed, err := parseStructFromPkg(p.pkg, structName)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, fmt.Errorf("struct not found: %q", structName)
	}

	parsed.Methods = parseMethods(p.pkg, structName)
	return parsed, nil
}

func parseStructFromPkg(pkg *packages.Package, structName string) (*generator.ParsedStruct, error) {
	for _, file := range pkg.Syntax {
		parsed, err := parseStructFromFile(file, structName)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			parsed.Package = generator.NewPackage(pkg.Types.Path(), pkg.Types.Name())
			return parsed, nil
		}
	}

	return nil, nil
}

func parseStructFromFile(file *ast.File, structName string) (*generator.ParsedStruct, error) {
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
			if structName != typeSpec.Name.Name {
				continue
			}

			switch typeSpec.Type.(type) {
			case *ast.StructType:
				t := typeSpec.Type.(*ast.StructType)
				return parseStruct(t, structName, file)
			default:
				return nil, fmt.Errorf(
					"specified declaration is not a struct: %T", typeSpec.Type,
				)
			}
		}
	}

	return nil, nil
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

func parseMethods(pkg *packages.Package, structName string) generator.Methods {
	result := generator.Methods{}
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			funcDecl, isFunc := d.(*ast.FuncDecl)
			if !isFunc {
				continue
			}
			if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
				continue
			}
			recv := funcDecl.Recv.List[0].Type
			if star, isStar := recv.(*ast.StarExpr); isStar {
				recv = star.X
			}
			ident, ok := recv.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name != structName {
				continue
			}
			result = append(result, generator.NewMethod(funcDecl.Name.Name))
		}
	}
	return result
}
