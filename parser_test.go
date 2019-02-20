package xml

import (
	"reflect"
	"testing"
)

func TestParser_Test(t *testing.T) {
	type fields struct {
		source string
		cursor int
	}
	tests := []struct {
		name   string
		fields fields
		str    string
		want   bool
	}{
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "a",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "ab",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "abc",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 1,
			},
			str:  "a",
			want: false,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 1,
			},
			str:  "bc",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 2,
			},
			str:  "c",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 2,
			},
			str:  "a",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := p.Tests(tt.str); got != tt.want {
				t.Errorf("Parser.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Get(t *testing.T) {
	type fields struct {
		source string
		cursor int
	}
	tests := []struct {
		name   string
		fields fields
		want   rune
	}{
		{
			fields: fields{
				source: "a",
				cursor: 0,
			},
			want: 'a',
		},
		{
			fields: fields{
				source: "a",
				cursor: 1,
			},
			want: EOF,
		},
		{
			fields: fields{
				source: "a",
				cursor: 2,
			},
			want: EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := p.Get(); got != tt.want {
				t.Errorf("Parser.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseXmlDecl(t *testing.T) {

	tests := []struct {
		name    string
		str     string
		want    *XMLDecl
		wantErr bool
	}{
		{
			str: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
			want: &XMLDecl{
				Version:    "1.0",
				Encoding:   "UTF-8",
				Standalone: true,
			},
		},
		{
			str: `<?xml version="2.0"  standalone="no"    ?>`,
			want: &XMLDecl{
				Version:    "2.0",
				Encoding:   "",
				Standalone: false,
			},
		},
		{
			str:     `<xml version="2.0" >`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.str)
			got, err := p.parseXmlDecl()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseXmlDecl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseXmlDecl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newRune(r rune) *rune { return &r }

func TestParser_parseProlog(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Prolog
		wantErr bool
	}{
		{
			name: "nesting elements",
			source: `
			<?xml version="1.0"?>
			<!DOCTYPE student [
				<!ELEMENT student (surname,firstname*,dob?,(origin|sex)?)>
				<!ELEMENT surname (#PCDATA)>
				<!ELEMENT firstname (#PCDATA)>
				<!ELEMENT sex (#PCDATA)>
			]>`,
			want: &Prolog{
				XMLDecl: &XMLDecl{
					Version: "1.0",
				},
				DOCType: &DOCType{
					Name: "student",
					Markups: []Markup{
						&Element{
							Name: "student",
							ContentSpec: &Children{
								ChoiceSeq: &Seq{
									CPs: []CP{
										CP{
											Name: "surname",
										},
										CP{
											Name:   "firstname",
											Suffix: newRune('*'),
										},
										CP{
											Name:   "dob",
											Suffix: newRune('?'),
										},
										CP{
											ChoiceSeq: &Choice{
												CPs: []CP{
													CP{
														Name: "origin",
													},
													CP{
														Name: "sex",
													},
												},
											},
											Suffix: newRune('?'),
										},
									},
								},
							},
						},
						&Element{
							Name:        "surname",
							ContentSpec: &Mixed{},
						},
						&Element{
							Name:        "firstname",
							ContentSpec: &Mixed{},
						},
						&Element{
							Name:        "sex",
							ContentSpec: &Mixed{},
						},
					},
				},
			},
		},
		{
			source: `<?xml version="1.0" standalone="yes" ?>

			<!--open the DOCTYPE declaration -
			  the open square bracket indicates an internal DTD-->
			<!DOCTYPE foo [
			
			<!--define the internal DTD-->
			  <!ELEMENT foo (#PCDATA)>
			
			<!--close the DOCTYPE declaration-->
			]>
			`,
			want: &Prolog{
				XMLDecl: &XMLDecl{
					Version:    "1.0",
					Standalone: true,
				},
				DOCType: &DOCType{
					Name: "foo",
					Markups: []Markup{
						Comment(`define the internal DTD`),
						&Element{
							Name:        "foo",
							ContentSpec: &Mixed{},
						},
						Comment(`close the DOCTYPE declaration`),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseProlog()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseProlog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseProlog() = %v, want %v", got, tt.want)
			}
		})
	}
}
