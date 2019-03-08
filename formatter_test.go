package xml

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

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
				},
				depth: 0,
			},
			want: `<!DOCTYPE html PUBLIC "public" "system" [
 <!ELEMENT code (#PCDATA)>
 <!NOTATION vrml PUBLIC "VRML 1.0">
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
			fmt.Println(buf.String(), tt.want)
			if buf.String() != tt.want {
				t.Errorf("want %q, but got %q", tt.want, buf.String())
			}
		})
	}
}
