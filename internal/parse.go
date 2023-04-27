package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func Parse(pathToFile, typeName string) (*ParsedStruct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, pathToFile, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, node := range f.Decls {
		genDecl, ok := node.(*ast.GenDecl)
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
					"specified declaration is of unsupported type %T, expected a struct", typeSpec.Type,
				)
			}
		}
	}

	return nil, fmt.Errorf("type %q not found", typeName)
}

func parseStruct(structType *ast.StructType, typeName string) (*ParsedStruct, error) {
	result := &ParsedStruct{
		Name:   typeName,
		Fields: ParsedFields{},
	}

	for _, field := range structType.Fields.List {
		var fieldType string
		switch field.Type.(type) {
		case *ast.Ident:
			fieldType = field.Type.(*ast.Ident).Name
		case *ast.StarExpr:
			fieldType = "*" + field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
		default:
			return nil, fmt.Errorf("field %q is of unknown type %T", field.Names, field.Type)
		}

		for _, name := range field.Names {
			result.Fields = append(result.Fields, ParsedField{
				Name: FieldName(name.Name),
				Type: FieldType(fieldType),
			})
		}
	}

	return result, nil
}

type ParsedStruct struct {
	Name   string
	Fields ParsedFields
}

type ParsedFields []ParsedField

type ParsedField struct {
	Name FieldName
	Type FieldType
}
