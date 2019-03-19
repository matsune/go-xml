package xml

import (
	"fmt"
	"io"
	"os"
)

func Format(a AST) {
	NewFormatter().Format(a)
}

type (
	Formatter struct {
		Indent string
		io.Writer
	}
)

func NewFormatter(opts ...fmtOption) *Formatter {
	f := &Formatter{
		Indent: "\t",
		Writer: os.Stdout,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type fmtOption func(*Formatter) error

func Indent(i string) fmtOption {
	return func(f *Formatter) error {
		f.Indent = i
		return nil
	}
}

func Writer(w io.Writer) fmtOption {
	return func(f *Formatter) error {
		f.Writer = w
		return nil
	}
}

func (f *Formatter) print(a ...interface{}) {
	fmt.Fprint(f.Writer, a...)
}

func (f *Formatter) printf(format string, a ...interface{}) {
	fmt.Fprintf(f.Writer, format, a...)
}

func (f *Formatter) println(a ...interface{}) {
	fmt.Fprintln(f.Writer, a...)
}

func (f *Formatter) ln() {
	fmt.Fprintln(f.Writer)
}

func (f *Formatter) insertIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Fprintf(f.Writer, "%s", f.Indent)
	}
}

func (f *Formatter) Format(a AST) {
	f.format(a, 0)
}

func (f *Formatter) FormatDepth(a AST, depth int) {
	f.format(a, depth)
}

func (f *Formatter) format(a AST, depth int) {
	if a == nil {
		return
	}
	switch v := a.(type) {
	case *XML:
		f.formatXML(v, depth)
	case *Prolog:
		f.formatProlog(v, depth)
	case *XMLDecl:
		f.formatXMLDecl(v, depth)
	case *DOCType:
		f.FormatDOCType(v, depth)
	case *Element:
		f.FormatElement(v, depth)
	case Terminal:
		f.insertIndent(depth)
		f.print(v.ToString())
	default:
		panic("unknown AST type")
	}
}

func (f *Formatter) formatXML(x *XML, depth int) {
	if x == nil {
		return
	}
	f.format(x.Prolog, depth)
	f.ln()
	f.format(x.Element, depth)
	f.ln()
	for _, m := range x.Misc {
		f.insertIndent(depth)
		f.format(m, depth)
	}
}

func (f *Formatter) formatProlog(p *Prolog, depth int) {
	if p == nil {
		return
	}
	f.format(p.XMLDecl, depth)
	f.ln()
	f.format(p.DOCType, depth)
}

func (f *Formatter) formatXMLDecl(x *XMLDecl, depth int) {
	if x == nil {
		return
	}
	f.insertIndent(depth)
	f.printf(`<?xml version="%s"`, x.Version)
	if len(x.Encoding) > 0 {
		f.printf(` encoding="%s"`, x.Encoding)
	}
	stdStr := "no"
	if x.Standalone {
		stdStr = "yes"
	}
	f.printf(" standalone=\"%s\" ?>", stdStr)
}

func (f *Formatter) FormatDOCType(d *DOCType, depth int) {
	if d == nil {
		return
	}
	f.insertIndent(depth)
	f.printf("<!DOCTYPE %s", d.Name)
	if d.ExtID != nil {
		f.printf(" %s", d.ExtID.ToString())
	}

	hasMarkups := len(d.Markups) > 0 || d.PERef != nil

	if hasMarkups {
		f.println(" [")
	}
	for _, m := range d.Markups {
		f.format(m, depth+1)
		f.ln()
	}
	if d.PERef != nil {
		f.format(d.PERef, depth+1)
		f.ln()
	}
	if hasMarkups {
		f.insertIndent(depth)
		f.print("]")
	}
	f.print(">")
}

func (f *Formatter) FormatElement(e *Element, depth int) {
	if e == nil {
		return
	}
	f.insertIndent(depth)
	f.printf("<%s", e.Name)

	for _, attr := range e.Attrs {
		f.printf(" %s", attr.ToString())
	}

	if e.IsEmptyTag {
		f.printf("/>")
		return
	}

	f.print(">")

	for i, c := range e.Contents {
		switch v := c.(type) {
		case *Element:
			f.ln()
			f.format(v, depth+1)
			if i == len(e.Contents)-1 {
				f.ln()
				f.insertIndent(depth)
			}
		case AST:
			f.ln()
			f.format(v, depth+1)
		default: // should be string only
			f.print(v)
		}
	}

	f.printf("</%s>", e.Name)
}
