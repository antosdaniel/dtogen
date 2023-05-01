package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"github.com/antosdaniel/dtogen/internal/generator"
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
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedFiles,
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
	result := &generator.ParsedStruct{
		Name:    typeName,
		Imports: fileImports(file),
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
			if !isMethod(funcDecl) {
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

func (p *parser) GetFunctions() (generator.ParsedFunctions, generator.Imports, error) {
	if p.pkg == nil {
		return nil, nil, fmt.Errorf("load package first")
	}

	funcs := generator.ParsedFunctions{}
	imports := generator.Imports{}
	for _, f := range p.pkg.Syntax {
		for _, d := range f.Decls {
			funcDecl, isFunc := d.(*ast.FuncDecl)
			if !isFunc {
				continue
			}
			if isMethod(funcDecl) {
				continue
			}
			// TODO: it's entire function, not just body
			funcBody, err := getSource(p.pkg, funcDecl.Pos(), funcDecl.Body.End())
			if err != nil {
				return nil, nil, err
			}
			funcs = append(funcs, generator.NewParsedFunction(funcDecl.Name.Name, funcBody))
			imports = append(imports, fileImports(f)...)
		}
	}
	return funcs, imports, nil
}

func isMethod(f *ast.FuncDecl) bool {
	return f.Recv != nil && len(f.Recv.List) > 0
}

func getSource(pkg *packages.Package, start, end token.Pos) (string, error) {
	startPos := pkg.Fset.Position(start)
	endPos := pkg.Fset.Position(end)
	if startPos.Filename != endPos.Filename {
		return "", fmt.Errorf("cant get source spanning multiple files")
	}

	file, err := getFile(startPos.Filename)
	if err != nil {
		return "", err
	}
	return file[startPos.Offset:endPos.Offset], nil
}
func getFile(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to load file %q: %w", filename, err)
	}

	return string(b), nil
}

func fileImports(file *ast.File) generator.Imports {
	imports := make(generator.Imports, len(file.Imports))
	for _, i := range file.Imports {
		imports = append(imports, generator.NewImport(i))
	}
	return imports
}
