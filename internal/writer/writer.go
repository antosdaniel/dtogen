package writer

import (
	"fmt"
	"go/ast"
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
func (w *Writer) In() *Writer {
	w.currentIndent += 1
	return w
}

// Out Decreases indentation. Next lines will be indented by one level (tab) less.
func (w *Writer) Out() *Writer {
	if w.currentIndent > 0 {
		w.currentIndent -= 1
	}
	return w
}

func (w *Writer) indent() string {
	indent := ""
	for i := 0; i < w.currentIndent; i++ {
		indent += "\t"
	}
	return indent
}

func (w *Writer) Write(code string) *Writer {
	w.body += w.indent() + code
	return w
}

func (w *Writer) WriteLine(code string) *Writer {
	if code == "" {
		return w.WriteEmptyLine()
	}
	w.body += w.indent() + code + "\n"
	return w
}

func (w *Writer) WriteEmptyLine() *Writer {
	w.body += "\n"
	return w
}

func (w *Writer) WritePackage(pkg string) *Writer {
	return w.WriteLine(fmt.Sprintf("package %s", pkg))
}

func (w *Writer) WriteStruct(s Struct) *Writer {
	w.WriteLine(fmt.Sprintf("type %s struct {", s.Name))
	w.In()
	col := typeShouldBeWrittenAtColumn(s.Fields)
	for _, f := range s.Fields {
		fieldType := writeType(f.Type)
		space := strings.Repeat(" ", col-len(f.Name))
		w.WriteLine(f.Name + space + fieldType)
	}
	w.Out()
	w.WriteLine("}")
	return w
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

func typeShouldBeWrittenAtColumn(fields Fields) int {
	result := 0
	for _, f := range fields {
		l := len(f.Name)
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
		return "*" + fieldType.(*ast.StarExpr).X.(*ast.Ident).Name
	default:
		panic(fmt.Errorf("unsupported type: %T", fieldType))
	}
}
