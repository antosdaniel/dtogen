package writer

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/antosdaniel/mappergen/internal/generator"
	"golang.org/x/tools/imports"
)

// writer Helps writing Go code.
type writer struct {
	body string

	currentIndent int
}

func New() generator.Writer {
	return &writer{
		body:          "",
		currentIndent: 0,
	}
}

// String Returns written code.
func (w *writer) String() string {
	result, err := imports.Process("", []byte(w.body), nil)
	if err == nil {
		return string(result)
	}
	return w.body
}

// In Increases indentation. Next lines will be indented by one level (tab) more.
func (w *writer) In() {
	w.currentIndent += 1
}

// Out Decreases indentation. Next lines will be indented by one level (tab) less.
func (w *writer) Out() {
	if w.currentIndent > 0 {
		w.currentIndent -= 1
	}
}

func (w *writer) indent() string {
	indent := ""
	for i := 0; i < w.currentIndent; i++ {
		indent += "\t"
	}
	return indent
}

func (w *writer) Write(code string) {
	w.body += w.indent() + code
}

func (w *writer) WriteLine(code string) {
	if code == "" {
		w.WriteEmptyLine()
		return
	}
	w.body += w.indent() + code + "\n"
}

func (w *writer) WriteEmptyLine() {
	w.body += "\n"
}

func (w *writer) WritePackage(pkg string) {
	w.WriteLine(fmt.Sprintf("package %s", pkg))
	w.WriteEmptyLine()
}

func (w *writer) WriteImports(imports generator.Imports) {
	if len(imports) == 0 {
		return
	}

	w.WriteLine("import (")
	w.In()
	for _, i := range imports {
		if i.Alias() != "" {
			w.WriteLine(fmt.Sprintf("%s %q", i.Alias(), i.Path()))
		} else {
			w.WriteLine(fmt.Sprintf("%q", i.Path()))
		}
	}
	w.Out()
	w.WriteLine(")")
	w.WriteEmptyLine()
}

func writeType(fieldType ast.Expr) string {
	switch fieldType.(type) {
	case *ast.Ident:
		return fieldType.(*ast.Ident).Name
	case *ast.StarExpr:
		return "*" + writeType(fieldType.(*ast.StarExpr).X)
	case *ast.SelectorExpr:
		sel := fieldType.(*ast.SelectorExpr)
		return writeType(sel.X) + "." + sel.Sel.Name
	case *ast.ArrayType:
		return "[]" + writeType(fieldType.(*ast.ArrayType).Elt)
	case *ast.MapType:
		m := fieldType.(*ast.MapType)
		return "map[" + writeType(m.Key) + "]" + writeType(m.Value)
	default:
		return "<unsupported type>"
	}
}

func (w *writer) WriteMapper(mapper generator.Mapper, outputPkg string) {
	funcName := "To" + strings.ToUpper(mapper.Dst().TypeName()[:1]) + mapper.Dst().TypeName()[1:]
	funcArgs := mapperArgs(mapper, outputPkg)
	returnType := structType(mapper.Dst(), outputPkg)
	w.WriteLine(fmt.Sprintf("func %s(%s) %s {", funcName, funcArgs, returnType))
	w.In()
	w.WriteLine(fmt.Sprintf("return %s{", returnType))
	w.In()
	for _, m := range mapper.Mappings() {
		src := mappingSrc(m.Src())
		w.WriteLine(fmt.Sprintf("%s: src.%s,", m.DstField(), src))
	}
	w.Out()
	w.WriteLine("}")
	w.Out()
	w.WriteLine("}")
}

func mapperArgs(mapper generator.Mapper, outputPkg string) string {
	src := mapper.Src()
	if len(src) <= 1 {
		return fmt.Sprintf("src %s", structType(src[0], outputPkg))
	}

	return "" // TODO: handle later
}

func structType(ms generator.MappedStruct, outputPkg string) string {
	if outputPkg == ms.Pkg().Name() {
		return ms.TypeName()
	}
	return fmt.Sprintf("%s.%s", ms.Pkg().Name(), ms.TypeName())
}

func mappingSrc(src generator.MappingSrc) string {
	switch src.(type) {
	case generator.MappingSrcField:
		return src.(generator.MappingSrcField).FieldName()
	default:
		return "<unsupported>"
	}
}
