package writer

import (
	"fmt"
	"github.com/antosdaniel/dtogen/internal/generator"
	"go/ast"
	"sort"
	"strings"
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
}

func (w *writer) WriteImports(imports generator.Imports) {
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

func (w *writer) WriteStruct(structName string, fields generator.Fields) {
	w.WriteLine(fmt.Sprintf("type %s struct {", structName))
	w.In()
	col := longestFieldNameLength(fields) + 1
	for _, f := range fields {
		fieldType := f.OverrideTypeTo()
		if fieldType == "" {
			fieldType = writeType(f.Type())
		}
		space := strings.Repeat(" ", col-len(f.DesiredName()))
		w.WriteLine(f.DesiredName() + space + fieldType)
	}
	w.Out()
	w.WriteLine("}")
}

func longestFieldNameLength(fields generator.Fields) int {
	result := 0
	for _, f := range fields {
		l := len(f.DesiredName())
		if l > result {
			result = l
		}
	}

	return result
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

func (w *writer) WriteMapper(mapper generator.Mapper) {
	w.WriteLine(fmt.Sprintf(
		"func New%s(src %s.%s) %s {",
		mapper.DestinationTypeName,
		mapper.SourceImportName,
		mapper.SourceTypeName,
		mapper.DestinationTypeName,
	))
	w.In()
	w.WriteLine(fmt.Sprintf("return %s{", mapper.DestinationTypeName))
	w.In()
	col := longestMappingDestinationLength(mapper.Mappings) + 1
	for _, m := range mapper.Mappings {
		space := strings.Repeat(" ", col-len(m.Destination()))
		src := m.Source()
		if m.IsMethod() {
			src += "()"
		}
		w.WriteLine(fmt.Sprintf("%s:%ssrc.%s,", m.Destination(), space, src))
	}
	w.Out()
	w.WriteLine("}")
	w.Out()
	w.WriteLine("}")
}

func longestMappingDestinationLength(mappings generator.Mappings) int {
	result := 0
	for _, m := range mappings {
		l := len(m.Destination())
		if l > result {
			result = l
		}
	}

	return result
}
