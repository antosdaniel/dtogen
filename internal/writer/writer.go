package writer

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"go/ast"
	"sort"
	"strings"
)

// Writer Helps writing Go code.
type Writer struct {
	body string

	currentIndent int
}

func New() *Writer {
	return &Writer{
		body:          "",
		currentIndent: 0,
	}
}

// String Returns written code.
func (w *Writer) String() string {
	return w.body
}

// In Increases indentation. Next lines will be indented by one level (tab) more.
func (w *Writer) In() {
	w.currentIndent += 1
}

// Out Decreases indentation. Next lines will be indented by one level (tab) less.
func (w *Writer) Out() {
	if w.currentIndent > 0 {
		w.currentIndent -= 1
	}
}

func (w *Writer) indent() string {
	indent := ""
	for i := 0; i < w.currentIndent; i++ {
		indent += "\t"
	}
	return indent
}

func (w *Writer) Write(code string) {
	w.body += w.indent() + code
}

func (w *Writer) WriteLine(code string) {
	if code == "" {
		w.WriteEmptyLine()
		return
	}
	w.body += w.indent() + code + "\n"
}

func (w *Writer) WriteEmptyLine() {
	w.body += "\n"
}

func (w *Writer) WritePackage(pkg string) {
	w.WriteLine(fmt.Sprintf("package %s", pkg))
}

func (w *Writer) WriteImports(imports generator.Imports) {
	if len(imports) == 0 {
		return
	}

	w.WriteLine("import (")
	w.In()
	for _, i := range sortImports(imports) {
		if i.Alias() != "" {
			w.WriteLine(fmt.Sprintf("%s %q", i.Alias(), i.Path()))
		} else {
			w.WriteLine(fmt.Sprintf("%q", i.Path()))
		}
	}
	w.Out()
	w.WriteLine(")")
}

func sortImports(imports generator.Imports) generator.Imports {
	sort.Slice(imports, func(i, j int) bool {
		// First sort by path A-Z
		if imports[i].Path() != imports[j].Path() {
			return imports[i].Path() < imports[j].Path()
		}

		// Then sort by aliases A-Z
		return imports[i].Alias() < imports[j].Alias()
	})

	return imports
}

func (w *Writer) WriteStruct(structName string, fields generator.Fields) {
	w.WriteLine(fmt.Sprintf("type %s struct {", structName))
	w.In()
	col := typeShouldBeWrittenAtColumn(fields)
	for _, f := range fields {
		fieldType := writeType(f.Type())
		space := strings.Repeat(" ", col-len(f.DesiredName()))
		w.WriteLine(f.DesiredName() + space + fieldType)
	}
	w.Out()
	w.WriteLine("}")
}

func typeShouldBeWrittenAtColumn(fields generator.Fields) int {
	result := 0
	for _, f := range fields {
		l := len(f.DesiredName())
		if l > result {
			result = l
		}
	}

	return result + 1
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
	default:
		panic(fmt.Errorf("unsupported type: %T", fieldType))
	}
}
