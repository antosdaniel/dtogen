package parser

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
)

type parser struct {
	pkg  *packages.Package
	file *ast.File
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

func (p *parser) LoadFile(filePath string) (generator.Parser, error) {
	fileSet := token.NewFileSet()
	file, err := goparser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	return &parser{file: file}, nil
}

// TODO: Refactor to FileParser and PackageParser, so we don't have to check both.

func (p *parser) GetStruct(typeName string) (*generator.ParsedStruct, error) {
	if p.pkg != nil {
		parsed, err := parseStructFromPkg(p.pkg, typeName)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			return parsed, nil
		}

		return nil, fmt.Errorf("type not found: %q", typeName)
	}
	if p.file != nil {
		parsed, err := parseStructFromFile(p.file, typeName)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			return parsed, nil
		}

		return nil, fmt.Errorf("type not found: %q", typeName)
	}

	return nil, fmt.Errorf("package nor file loaded")
}

func parseStructFromPkg(pkg *packages.Package, typeName string) (*generator.ParsedStruct, error) {
	for _, file := range pkg.Syntax {
		parsed, err := parseStructFromFile(file, typeName)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			return parsed, nil
		}
	}

	return nil, nil
}

func parseStructFromFile(file *ast.File, typeName string) (*generator.ParsedStruct, error) {
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
