package xml

import (
	"fmt"
	"io"
	"os"
)

type (
	Formatter struct {
		Indent string
		io.Writer
	}
)

func NewFormatter() *Formatter {
	return &Formatter{
		Indent: "\t",
		Writer: os.Stdout,
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

func (f *Formatter) FormatXML(x *XML) {
	if x == nil {
		return
	}
	f.FormatProlog(x.Prolog, 0)
	f.ln()
	f.FormatElement(x.Element, 0)
	f.ln()
	for _, m := range x.Misc {
		f.println(m)
	}
}

func (f *Formatter) FormatProlog(p *Prolog, depth int) {
	if p == nil {
		return
	}
	f.FormatXMLDecl(p.XMLDecl, depth)
	f.ln()
	f.FormatDOCType(p.DOCType, depth)
}

func (f *Formatter) FormatXMLDecl(x *XMLDecl, depth int) {
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
		f.printf(" %s", d.ExtID)
	}

	hasMarkups := len(d.Markups) > 0 || d.PERef != nil

	if hasMarkups {
		f.print(" [\n")
	}
	for _, m := range d.Markups {
		f.insertIndent(depth + 1)
		f.println(m)
	}
	if d.PERef != nil {
		f.insertIndent(depth + 1)
		f.println(d.PERef)
	}
	if hasMarkups {
		f.insertIndent(depth)
		f.print("]")
	}
	f.printf(">")
}

func (f *Formatter) FormatElement(e *Element, depth int) {
	if e == nil {
		return
	}
	f.insertIndent(depth)
	f.printf("<%s", e.Name)

	for _, attr := range e.Attrs {
		f.printf(" %s", attr)
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
			f.FormatElement(v, depth+1)
			if i == len(e.Contents)-1 {
				f.ln()
				f.insertIndent(depth)
			}
		default:
			f.print(v, depth+1)
		}
	}

	f.printf("</%s>", e.Name)
}
