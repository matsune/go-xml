package xml

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestIndent(t *testing.T) {
	tests := []struct {
		name string
		i    string
	}{
		{
			name: "tab",
			i:    "\t",
		},
		{
			name: "space",
			i:    "  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(Indent(tt.i))
			if f.Indent != tt.i {
				t.Errorf("Indent = %v, want %v", f.Indent, tt.i)
			}
		})
	}
}

func TestWriter(t *testing.T) {
	tests := []struct {
		name string
		w    io.Writer
	}{
		{
			name: "stdOut",
			w:    os.Stdout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(Writer(tt.w))
			if f.Writer != tt.w {
				t.Errorf("Indent = %v, want %v", f.Writer, tt.w)
			}
		})
	}
}

func TestFormatter_FormatXMLDecl(t *testing.T) {
	type args struct {
		x     *XMLDecl
		depth int
	}
	tests := []struct {
		name   string
		Indent string
		args   args
		want   string
	}{
		{
			name: "nil",
			args: args{
				x: nil,
			},
			want: "",
		},
		{
			Indent: " ",
			args: args{
				x: &XMLDecl{
					Version:    "1.0",
					Encoding:   "UTF-8",
					Standalone: true,
				},
				depth: 2,
			},
			want: `  <?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			f := &Formatter{
				Indent: tt.Indent,
				Writer: w,
			}
			f.formatXMLDecl(tt.args.x, tt.args.depth)
			w.Flush()
			if buf.String() != tt.want {
				t.Errorf("want %q, but got %q", tt.want, buf.String())
			}
		})
	}
}

func TestFormatter_FormatDOCType(t *testing.T) {
	type args struct {
		d     *DOCType
		depth int
	}
	tests := []struct {
		name   string
		Indent string
		args   args
		want   string
	}{
		{
			name: "nil",
			args: args{
				d: nil,
			},
			want: "",
		},
		{
			Indent: " ",
			args: args{
				d: &DOCType{
					Name: "html",
					ExtID: &ExternalID{
						Type:   ExternalTypePublic,
						Pubid:  "public",
						System: "system",
					},
					Markups: []Markup{
						&ElementDecl{
							Name:        "code",
							ContentSpec: &Mixed{},
						},
						&Notation{
							Name: "vrml",
							ExtID: ExternalID{
								Type:  ExternalTypePublic,
								Pubid: "VRML 1.0",
							},
						},
					},
					PERef: &PERef{
						Name: "a",
					},
				},
				depth: 0,
			},
			want: `<!DOCTYPE html PUBLIC "public" "system" [
 <!ELEMENT code (#PCDATA)>
 <!NOTATION vrml PUBLIC "VRML 1.0">
 %a;
]>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			f := &Formatter{
				Indent: tt.Indent,
				Writer: w,
			}
			f.FormatDOCType(tt.args.d, tt.args.depth)
			w.Flush()
			if buf.String() != tt.want {
				t.Errorf("want %q, but got %q", tt.want, buf.String())
			}
		})
	}
}

func TestFormatter_FormatElement(t *testing.T) {
	type args struct {
		e     *Element
		depth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "element nil",
			args: args{
				e: nil,
			},
			want: "",
		},
		{
			name: "empty tag <empty/>",
			args: args{
				e: &Element{
					Name:       "empty",
					IsEmptyTag: true,
				},
			},
			want: "<empty/>",
		},
		{
			name: "empty with attrs",
			args: args{
				e: &Element{
					Name: "empty",
					Attrs: Attributes{
						&Attribute{
							Name: "attr",
							AttValue: AttValue{
								"attvalue",
							},
						},
					},
					IsEmptyTag: true,
				},
			},
			want: `<empty attr="attvalue"/>`,
		},
		{
			name: "nested element",
			args: args{
				e: &Element{
					Name: "root",
					Contents: []interface{}{
						&Element{
							Name:       "child",
							IsEmptyTag: true,
						},
						&PERef{
							Name: "peref",
						},
						Comment("comment"),
						&Element{
							Name: "child2",
							Contents: []interface{}{
								"string",
							},
						},
					},
				},
			},
			want: `<root>
	<child/>
	%peref;
	<!--comment-->
	<child2>string</child2>
</root>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			f := &Formatter{
				Indent: "\t",
				Writer: w,
			}
			f.FormatElement(tt.args.e, tt.args.depth)
			w.Flush()
			if buf.String() != tt.want {
				t.Errorf("want %q, but got %q", tt.want, buf.String())
			}
		})
	}
}
