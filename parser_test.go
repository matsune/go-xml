package xml

import (
	"fmt"
	"reflect"
	"testing"
)

func newRune(r rune) *rune { return &r }

func TestParser_parseProlog(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Prolog
		wantErr bool
	}{
		{
			name:    "error while parsing XmlDecl",
			source:  `<?xml version=?> <!DOCTYPE html >`,
			wantErr: true,
		},
		{
			name:    "error while parsing doctype",
			source:  `<?xml version="1.0" ?>  <!DOCTYPE html `,
			wantErr: true,
		},
		{
			name: "xml decl and dtd decl",
			source: `<?xml version="1.0"?>
			<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
			 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
			want: &Prolog{
				XMLDecl: &XMLDecl{
					Version: "1.0",
				},
				DOCType: &DOCType{
					Name: "html",
					ExtID: &ExternalID{
						Identifier: ExtPublic,
						Pubid:      "-//W3C//DTD XHTML 1.0 Transitional//EN",
						System:     "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd",
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

func TestParser_parseXmlDecl(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *XMLDecl
		wantErr bool
	}{
		{
			name:    "not <?xml",
			str:     `<xml version="1.0" standalone="no" ?>`,
			wantErr: true,
		},
		{
			name:    "error version",
			str:     `<?xml version=1.0 standalone="no" ?>`,
			wantErr: true,
		},
		{
			name:    "error encoding",
			str:     `<?xml version="1.0" encoding= ?>`,
			wantErr: true,
		},
		{
			name:    "error standalone",
			str:     `<?xml version="1.0" standalone= ?>`,
			wantErr: true,
		},
		{
			name:    "error close",
			str:     `<?xml version="1.0"  >`,
			wantErr: true,
		},
		{
			str: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
			want: &XMLDecl{
				Version:    "1.0",
				Encoding:   "UTF-8",
				Standalone: true,
			},
		},
		{
			str: `<?xml version="1.1"  standalone="no"    ?>`,
			want: &XMLDecl{
				Version:    "1.1",
				Encoding:   "",
				Standalone: false,
			},
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

func TestParser_parseVersion(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantVer string
		wantErr bool
	}{
		{
			name:    "not starts with spaces",
			source:  `version="1.0"`,
			wantErr: true,
		},
		{
			name:    "not starts with version",
			source:  ` ver="1.0"`,
			wantErr: true,
		},
		{
			name:    "not equal",
			source:  ` version:"1.0"`,
			wantErr: true,
		},
		{
			name:    "no quote",
			source:  ` version=1.0`,
			wantErr: true,
		},
		{
			name:    "error while parsing version num",
			source:  ` version=""`,
			wantErr: true,
		},
		{
			name:    "different quotes",
			source:  ` version="1.0'`,
			wantErr: true,
		},
		{
			source:  ` version="1.0" `,
			wantVer: "1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			gotVer, err := p.parseVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVer != tt.wantVer {
				t.Errorf("Parser.parseVersion() = %v, want %v", gotVer, tt.wantVer)
			}
		})
	}
}

func TestParser_parseMisc(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "error while parsing comment",
			source:  `<!--- comment`,
			wantErr: true,
		},
		{
			name:    "error while parsing PI",
			source:  `<? PI`,
			wantErr: true,
		},
		{
			name:    "starts with unknown keyword",
			source:  `< aa >`,
			wantErr: true,
		},
		{
			name:   "parse comment",
			source: `<!-- comment -->`,
			want:   Comment(" comment "),
		},
		{
			name:   "parse PI",
			source: `<?target aaa?>`,
			want: &PI{
				Target:      "target",
				Instruction: "aaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseMisc()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseMisc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseMisc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseName(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "starts with null char",
			source:  fmt.Sprintf("%caaa", 0),
			wantErr: true,
		},
		{
			source: ":abc",
			want:   ":abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseName()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEntityValue(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    EntityValue
		wantErr bool
	}{
		{
			name:    "not starts with quote",
			source:  `%aa;"`,
			wantErr: true,
		},
		{
			name:    "error while parsing charRef",
			source:  `"&#a"`,
			wantErr: true,
		},
		{
			name:    "error while parsing PERef",
			source:  `"%a"`,
			wantErr: true,
		},
		{
			name:   "entity ref",
			source: `"&a;"`,
			want: EntityValue{
				&EntityRef{
					Name: "a",
				},
			},
		},
		{
			name:   "multiple refs",
			source: `"a%b;c&#xdd;"`,
			want: EntityValue{
				"a",
				&PERef{
					Name: "b",
				},
				"c",
				&CharRef{
					Prefix: "&#x",
					Value:  "dd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEntityValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEntityValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseEntityValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseAttValue(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    AttValue
		wantErr bool
	}{
		{
			name:   "empty",
			source: `""`,
			want:   AttValue{},
		},
		{
			name:    "not started with quote",
			source:  `a"`,
			wantErr: true,
		},
		{
			name:    "error single <",
			source:  `"<"`,
			wantErr: true,
		},
		{
			name:    "error single &",
			source:  `"&"`,
			wantErr: true,
		},
		{
			name:   "EntityRef",
			source: `"&a;"`,
			want: AttValue{
				&EntityRef{Name: "a"},
			},
		},
		{
			name:   "CharRef",
			source: `"&#11;"`,
			want: AttValue{
				&CharRef{Prefix: "&#", Value: "11"},
			},
		},
		{
			name:   "multiple values",
			source: `"bb&#11;&#x20;&a;"`,
			want: AttValue{
				"bb",
				&CharRef{Prefix: "&#", Value: "11"},
				&CharRef{Prefix: "&#x", Value: "20"},
				&EntityRef{Name: "a"},
			},
		},
		{
			source: `"abcd"`,
			want:   AttValue{"abcd"},
		},
		{
			name:    "different quotes",
			source:  `"abcd'`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseAttValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseAttValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseAttValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseSystemLiteral(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "not starts with quote",
			source:  `aaa`,
			wantErr: true,
		},
		{
			name:    "different quotes",
			source:  `'a"`,
			wantErr: true,
		},
		{
			name:   "empty literal",
			source: `''`,
			want:   "",
		},
		{
			source: `"aaa"`,
			want:   "aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseSystemLiteral()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSystemLiteral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseSystemLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parsePubidLiteral(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "not starts with quote",
			source:  `aaa`,
			wantErr: true,
		},
		{
			name:    "different quotes",
			source:  `'a"`,
			wantErr: true,
		},
		{
			name:   "empty literal",
			source: `''`,
			want:   "",
		},
		{
			source: `"aaa"`,
			want:   "aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parsePubidLiteral()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parsePubidLiteral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parsePubidLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseComment(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    Comment
		wantErr bool
	}{
		{
			name:    "not starts with <!--",
			source:  "<-- aa -->",
			wantErr: true,
		},
		{
			name:    "not end with -->",
			source:  "<!-- aa --",
			wantErr: true,
		},
		{
			name:    "contains not char",
			source:  fmt.Sprintf("<!-- %c -->", 0x0),
			wantErr: true,
		},
		{
			source: "<!-- this is comment-->",
			want:   Comment(" this is comment"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseComment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parsePI(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *PI
		wantErr bool
	}{
		{
			name:    "not starts with <?",
			source:  "<target?>",
			wantErr: true,
		},
		{
			name:    "error parsing target",
			source:  `<?" ?>`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `<?target `,
			wantErr: true,
		},
		{
			source: `<?target key="value" ?>`,
			want: &PI{
				Target:      "target",
				Instruction: `key="value" `,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parsePI()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parsePI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parsePI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parsePITarget(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "invalid name",
			source:  "",
			wantErr: true,
		},
		{
			name:    "contains x",
			source:  "wordx",
			wantErr: true,
		},
		{
			name:    "contains X",
			source:  "worXd",
			wantErr: true,
		},
		{
			name:    "contains m",
			source:  "womrd",
			wantErr: true,
		},
		{
			name:    "contains M",
			source:  "wMord",
			wantErr: true,
		},
		{
			name:    "contains l",
			source:  "lword",
			wantErr: true,
		},
		{
			name:    "contains l",
			source:  "wordL",
			wantErr: true,
		},
		{
			name:   "not contain xmlXML",
			source: "word",
			want:   "word",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parsePITarget()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parsePITarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parsePITarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseDoctype(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *DOCType
		wantErr bool
	}{
		{
			name:    "not starts with <!DOCTYPE",
			source:  `<DOC name >`,
			wantErr: true,
		},
		{
			name:    "no space between <!DOCTYPE and name",
			source:  `<!DOCTYPEname >`,
			wantErr: true,
		},
		{
			name:    "invalid name",
			source:  `<!DOCTYPE ' >`,
			wantErr: true,
		},
		{
			name:    "invalid externalID",
			source:  `<!DOCTYPE name PUBLIC >`,
			wantErr: true,
		},
		{
			name: "error while parsing markup",
			source: `<!DOCTYPE name [
				<!ELEMENT 
			] >`,
			wantErr: true,
		},
		{
			name: "error while parsing PERef",
			source: `<!DOCTYPE name [
				%PEREF
			] >`,
			wantErr: true,
		},
		{
			name:    "not closed ]",
			source:  `<!DOCTYPE name [`,
			wantErr: true,
		},
		{
			name:    "not closed tag",
			source:  `<!DOCTYPE name []`,
			wantErr: true,
		},
		{
			name: "author element and entity",
			source: `<!DOCTYPE author [
				<!ELEMENT author (#PCDATA)>
				<!ENTITY js "Jo Smith">
			  ]>`,
			want: &DOCType{
				Name: "author",
				Markups: []Markup{
					&Element{
						Name:        "author",
						ContentSpec: &Mixed{},
					},
					&Entity{
						Name: "js",
						Type: EntityType_GE,
						Value: EntityValue{
							"Jo Smith",
						},
					},
				},
			},
		},
		{
			name: "author element and entity with externalID",
			source: `<!DOCTYPE author SYSTEM 'system' [
				<!ELEMENT author (#PCDATA)>
				<!ENTITY js "Jo Smith">
			  ]>`,
			want: &DOCType{
				Name: "author",
				ExtID: &ExternalID{
					Identifier: ExtSystem,
					System:     "system",
				},
				Markups: []Markup{
					&Element{
						Name:        "author",
						ContentSpec: &Mixed{},
					},
					&Entity{
						Name: "js",
						Type: EntityType_GE,
						Value: EntityValue{
							"Jo Smith",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseDoctype()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseDoctype() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseDoctype() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseMarkup(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    Markup
		wantErr bool
	}{
		{
			name:    "unknown markup",
			source:  "",
			wantErr: true,
		},
		{
			name:   "element",
			source: "<!ELEMENT student (id|(a,b)?)>",
			want: &Element{
				Name: "student",
				ContentSpec: &Children{
					ChoiceSeq: &Choice{
						CPs: []CP{
							CP{
								Name: "id",
							},
							CP{
								ChoiceSeq: &Seq{
									CPs: []CP{
										CP{
											Name: "a",
										},
										CP{
											Name: "b",
										},
									},
								},
								Suffix: newRune('?'),
							},
						},
					},
				},
			},
		},
		{
			name:   "attlist",
			source: "<!ATTLIST task status (important|normal) #REQUIRED>",
			want: &Attlist{
				Name: "task",
				Defs: []*AttDef{
					&AttDef{
						Name: "status",
						Type: &Enum{
							Cases: []string{"important", "normal"},
						},
						Decl: &DefaultDecl{
							Type: REQUIRED,
						},
					},
				},
			},
		},
		{
			name: "entity",
			source: `<!ENTITY a SYSTEM
			"http://example.com/a.gif">`,
			want: &Entity{
				Name: "a",
				Type: EntityType_GE,
				ExtID: &ExternalID{
					Identifier: ExtSystem,
					System:     "http://example.com/a.gif",
				},
			},
		},
		{
			name:   "notation",
			source: `<!NOTATION vrml PUBLIC "VRML 1.0">`,
			want: &Notation{
				Name: "vrml",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "VRML 1.0",
				},
			},
		},
		{
			name:   "PI",
			source: `<?target function="enable"?>`,
			want: &PI{
				Target:      "target",
				Instruction: `function="enable"`,
			},
		},
		{
			name:   "Comment",
			source: `<!-- this is comment -->`,
			want:   Comment(" this is comment "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseMarkup()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseMarkup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseMarkup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseStandalone(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    bool
		wantErr bool
	}{
		{
			name:    "not start with spaces",
			source:  `standalone='yes'`,
			wantErr: true,
		},
		{
			name:    "no standalone",
			source:  ` stand='yes'`,
			wantErr: true,
		},
		{
			name:    "error parse =",
			source:  ` standalone:'yes'`,
			wantErr: true,
		},
		{
			name:    "no quote",
			source:  ` standalone=yes`,
			wantErr: true,
		},
		{
			name:    "invalid bool value",
			source:  ` standalone='true'`,
			wantErr: true,
		},
		{
			name:    "difference quotes",
			source:  ` standalone="yes'`,
			wantErr: true,
		},
		{
			source: ` standalone="yes"`,
			want:   true,
		},
		{
			source: ` standalone='no'`,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseStandalone()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseStandalone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseStandalone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseElement(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Element
		wantErr bool
	}{
		{
			name:    "invalid head",
			source:  `<ELEMENT student (id)>`,
			wantErr: true,
		},
		{
			name:    "no space",
			source:  `<!ELEMENTstudent (id)>`,
			wantErr: true,
		},
		{
			name:    "name error",
			source:  `<!ELEMENT (id)>`,
			wantErr: true,
		},
		{
			name:    "no space after name name",
			source:  `<!ELEMENT student(id)>`,
			wantErr: true,
		},
		{
			name:    "invalid content spec",
			source:  `<!ELEMENT student ()>`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `<!ELEMENT student (id)`,
			wantErr: true,
		},
		{
			name:   "simple element",
			source: `<!ELEMENT student (id)>`,
			want: &Element{
				Name: "student",
				ContentSpec: &Children{
					ChoiceSeq: &Choice{
						CPs: []CP{
							CP{
								Name: "id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseElement()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseContentSpec(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    ContentSpec
		wantErr bool
	}{
		{
			name:   "parse EMPTY",
			source: "EMPTY",
			want:   &EMPTY{},
		},
		{
			name:   "parse ANY",
			source: "ANY",
			want:   &ANY{},
		},
		{
			name:   "parse mixed",
			source: "(#PCDATA)",
			want:   &Mixed{},
		},
		{
			name:   "parse children",
			source: "(id)",
			want: &Children{
				ChoiceSeq: &Choice{
					CPs: []CP{
						CP{
							Name: "id",
						},
					},
				},
			},
		},
		{
			name:    "error",
			source:  "(id",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseContentSpec()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseContentSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseContentSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseChildren(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Children
		wantErr bool
	}{
		{
			name:    "parse error",
			source:  "id|name)",
			wantErr: true,
		},
		{
			name:   "choice children",
			source: "(id|name)",
			want: &Children{
				ChoiceSeq: &Choice{
					CPs: []CP{
						CP{
							Name: "id",
						},
						CP{
							Name: "name",
						},
					},
				},
			},
		},
		{
			name:   "choice children with suffix",
			source: "(id|name)+",
			want: &Children{
				ChoiceSeq: &Choice{
					CPs: []CP{
						CP{
							Name: "id",
						},
						CP{
							Name: "name",
						},
					},
				},
				Suffix: newRune('+'),
			},
		},
		{
			name:   "seq children with suffix",
			source: "(id,name)+",
			want: &Children{
				ChoiceSeq: &Seq{
					CPs: []CP{
						CP{
							Name: "id",
						},
						CP{
							Name: "name",
						},
					},
				},
				Suffix: newRune('+'),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseChildren()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseChildren() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseCP(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *CP
		wantErr bool
	}{
		{
			name:    "error while parsing seq",
			source:  "(surname,)",
			wantErr: true,
		},
		{
			name:   "nested",
			source: "(surname,(origin|sex)?)",
			want: &CP{
				ChoiceSeq: &Seq{
					CPs: []CP{
						CP{
							Name: "surname",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseCP()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseChoice(t *testing.T) {
	type fields struct {
		Scanner *Scanner
	}
	tests := []struct {
		name    string
		source  string
		want    *Choice
		wantErr bool
	}{
		{
			name:    "not starts with (",
			source:  `surname|firstname)`,
			wantErr: true,
		},
		{
			name:    "error no cp",
			source:  `()`,
			wantErr: true,
		},
		{
			name: "error parsing cp",
			source: `(surname|
				)`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `(surname|firstname`,
			wantErr: true,
		},
		{
			source: `(surname|firstname)`,
			want: &Choice{
				CPs: []CP{
					CP{
						Name: "surname",
					},
					CP{
						Name: "firstname",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseChoice()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseChoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseSeq(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Seq
		wantErr bool
	}{
		{
			name:    "not starts with (",
			source:  `surname,firstname*,dob?,(origin|sex)?)`,
			wantErr: true,
		},
		{
			name:    "error no cp",
			source:  `()`,
			wantErr: true,
		},
		{
			name: "error parsing cp",
			source: `(surname,
				)`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `(surname,firstname*`,
			wantErr: true,
		},
		{
			source: `(surname,firstname*)`,
			want: &Seq{
				CPs: []CP{
					CP{
						Name: "surname",
					},
					CP{
						Name:   "firstname",
						Suffix: newRune('*'),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseSeq()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSeq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseSeq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseMixed(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Mixed
		wantErr bool
	}{
		{
			name:    "not starts with (",
			source:  "#PCDATA)*",
			wantErr: true,
		},
		{
			name:    "no #PCDATA",
			source:  "()*",
			wantErr: true,
		},
		{
			name:    "not | between #PCDATA and )",
			source:  "(#PCDATA  a)*",
			wantErr: true,
		},
		{
			name:    "no name after #PCDATA|",
			source:  "(#PCDATA|)*",
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  "(#PCDATA|a",
			wantErr: true,
		},
		{
			name:    "error if has names but no *",
			source:  "(#PCDATA|a)",
			wantErr: true,
		},
		{
			source: "(#PCDATA|a  |  b)*",
			want: &Mixed{
				Names: []string{"a", "b"},
			},
		},
		{
			source: "(  #PCDATA  )",
			want:   &Mixed{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseMixed()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseMixed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseMixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseAttlist(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Attlist
		wantErr bool
	}{
		{
			name:    "not starts with <!ATTLIST",
			source:  `<ATTLIST n >`,
			wantErr: true,
		},
		{
			name:    "no space",
			source:  `<!ATTLISTname >`,
			wantErr: true,
		},
		{
			name: "error parsing name",
			source: `<!ATTLIST 
			 >`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `<!ATTLIST name `,
			wantErr: true,
		},
		{
			name:    "not closed with AttDef",
			source:  `<!ATTLIST name height CDATA #REQUIRED `,
			wantErr: true,
		},
		{
			name:    "error parsing AttDef",
			source:  `<!ATTLIST name ' >`,
			wantErr: true,
		},
		{
			name:   "no attdefs",
			source: `<!ATTLIST name >`,
			want: &Attlist{
				Name: "name",
			},
		},
		{
			name:   "with attdefs",
			source: `<!ATTLIST image height CDATA #REQUIRED>`,
			want: &Attlist{
				Name: "image",
				Defs: []*AttDef{
					&AttDef{
						Name: "height",
						Type: Att_CDATA,
						Decl: &DefaultDecl{
							Type: REQUIRED,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseAttlist()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseAttlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseAttlist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseAttDef(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *AttDef
		wantErr bool
	}{
		{
			name:    "no space",
			source:  `name CDATA #REQUIRED`,
			wantErr: true,
		},
		{
			name:    "error parsing name",
			source:  ` .. CDATA #REQUIRED`,
			wantErr: true,
		},
		{
			name:    "no space after name",
			source:  ` name'CDATA #REQUIRED`,
			wantErr: true,
		},
		{
			name:    "error parsing AttType",
			source:  ` name CDA #REQUIRED`,
			wantErr: true,
		},
		{
			name:    "no space after AttType",
			source:  ` name CDATA#REQUIRED`,
			wantErr: true,
		},
		{
			name:    "error parsing DefaultDecl",
			source:  ` name CDATA #required`,
			wantErr: true,
		},
		{
			source: ` name CDATA #REQUIRED`,
			want: &AttDef{
				Name: "name",
				Type: Att_CDATA,
				Decl: &DefaultDecl{
					Type: REQUIRED,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseAttDef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseAttDef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseAttDef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseAttType(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    AttType
		wantErr bool
	}{
		{
			source: "CDATA",
			want:   Att_CDATA,
		},
		{
			source: "ID",
			want:   Att_ID,
		},
		{
			source: "IDREF",
			want:   Att_IDREF,
		},
		{
			source: "IDREFS",
			want:   Att_IDREFS,
		},
		{
			source: "ENTITY",
			want:   Att_ENTITY,
		},
		{
			source: "ENTITIES",
			want:   Att_ENTITIES,
		},
		{
			source: "NMTOKEN",
			want:   Att_NMTOKEN,
		},
		{
			source: "NMTOKENS",
			want:   Att_NMTOKENS,
		},
		{
			source: "NOTATION (a)",
			want: &NotationType{
				Names: []string{"a"},
			},
		},
		{
			source: "(a|b)",
			want: &Enum{
				Cases: []string{"a", "b"},
			},
		},
		{
			name:    "error empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseAttType()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseAttType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseAttType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseNotationType(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *NotationType
		wantErr bool
	}{
		{
			name:    "not starts with NOTATION",
			source:  " (test)",
			wantErr: true,
		},
		{
			name:    "no space",
			source:  "NOTATION(test)",
			wantErr: true,
		},
		{
			name:    "no (",
			source:  "NOTATION test",
			wantErr: true,
		},
		{
			name:    "no name",
			source:  "NOTATION ()",
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  "NOTATION (test",
			wantErr: true,
		},
		{
			name:    "not separated |",
			source:  "NOTATION (a,b)",
			wantErr: true,
		},
		{
			name:    "no name",
			source:  "NOTATION (a|)",
			wantErr: true,
		},
		{
			source: "NOTATION (a|b)",
			want: &NotationType{
				Names: []string{"a", "b"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseNotationType()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseNotationType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseNotationType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEnum(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Enum
		wantErr bool
	}{
		{
			name:    "not starts with (",
			source:  "a|b",
			wantErr: true,
		},
		{
			name:    "no name",
			source:  "()",
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  "(a",
			wantErr: true,
		},
		{
			name:    "not separated |",
			source:  "(a,b)",
			wantErr: true,
		},
		{
			name:    "no name",
			source:  "(a|)",
			wantErr: true,
		},
		{
			source: "(important|normal)",
			want: &Enum{
				Cases: []string{"important", "normal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEnum()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEnum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseEnum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseDefaultDecl(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *DefaultDecl
		wantErr bool
	}{
		{
			source: "#REQUIRED",
			want: &DefaultDecl{
				Type: REQUIRED,
			},
		},
		{
			source: "#IMPLIED",
			want: &DefaultDecl{
				Type: IMPLIED,
			},
		},
		{
			source: `#FIXED "a"`,
			want: &DefaultDecl{
				Type:     FIXED,
				AttValue: []interface{}{"a"},
			},
		},
		{
			name:    "no space",
			source:  `#FIXED"a"`,
			wantErr: true,
		},
		{
			name:    "error AttValue",
			source:  `#FIXED aa`,
			wantErr: true,
		},
		{
			name:   "no #FIXED",
			source: `"aa"`,
			want: &DefaultDecl{
				Type:     FIXED,
				AttValue: []interface{}{"aa"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseDefaultDecl()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseDefaultDecl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseDefaultDecl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseCharRef(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *CharRef
		wantErr bool
	}{
		{
			name:    "not starts with &#",
			source:  `&a;`,
			wantErr: true,
		},
		{
			name:    "not number after &#",
			source:  `&#a;`,
			wantErr: true,
		},
		{
			name:    "not number or alphabet after &#x",
			source:  `&#x.;`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `&#xa`,
			wantErr: true,
		},
		{
			source: `&#xcc;`,
			want: &CharRef{
				Prefix: "&#x",
				Value:  "cc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseCharRef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseCharRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseCharRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEntityReference(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *EntityRef
		wantErr bool
	}{
		{
			name:    "not starts with &",
			source:  `name;`,
			wantErr: true,
		},
		{
			name:    "error parse name",
			source:  `&;`,
			wantErr: true,
		},
		{
			name:    "not closed ;",
			source:  `&name`,
			wantErr: true,
		},
		{
			source: `&name;`,
			want: &EntityRef{
				Name: "name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEntityRef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEntityReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseEntityReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parsePERef(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *PERef
		wantErr bool
	}{
		{
			name:    "not starts with %",
			source:  `name;`,
			wantErr: true,
		},
		{
			name:    "error parse name",
			source:  `%;`,
			wantErr: true,
		},
		{
			name:    "not closed ;",
			source:  `%name`,
			wantErr: true,
		},
		{
			source: `%name;`,
			want: &PERef{
				Name: "name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parsePERef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parsePERef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parsePERef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEntity(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Entity
		wantErr bool
	}{
		{
			name:    "not starts with <!ENTITY",
			source:  `<ENTITY a "b,c,d">`,
			wantErr: true,
		},
		{
			name:    "no spaces after <!ENTITY",
			source:  `<!ENTITYa "b,c,d">`,
			wantErr: true,
		},
		{
			name:    "no space between % and name",
			source:  `<!ENTITY %a  "b,c,d">`,
			wantErr: true,
		},
		{
			name:    "error while parsing name",
			source:  `<!ENTITY ..  "b,c,d">`,
			wantErr: true,
		},
		{
			name:    "no space after name",
			source:  `<!ENTITY a"b,c,d">`,
			wantErr: true,
		},
		{
			name:    "error while parsing PEDef",
			source:  `<!ENTITY % a ">`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `<!ENTITY % a "b,c"`,
			wantErr: true,
		},
		{
			name:   "PEDecl",
			source: `<!ENTITY % a "b,c">`,
			want: &Entity{
				Name:  "a",
				Type:  EntityType_PE,
				Value: EntityValue{"b,c"},
			},
		},
		{
			name:   "GEDecl",
			source: `<!ENTITY a "b,c">`,
			want: &Entity{
				Name:  "a",
				Type:  EntityType_GE,
				Value: EntityValue{"b,c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEntity()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEntityDef(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    EntityValue
		want1   *ExternalID
		want2   string
		wantErr bool
	}{
		{
			name:   "EntityValue",
			source: `'ab%aa;'`,
			want: EntityValue{
				"ab",
				&PERef{
					Name: "aa",
				},
			},
		},
		{
			name:   "ExternalID",
			source: `SYSTEM "aa"`,
			want1: &ExternalID{
				Identifier: ExtSystem,
				System:     "aa",
			},
		},
		{
			name:   "ExternalID with NData",
			source: `SYSTEM "aa" NDATA bb`,
			want1: &ExternalID{
				Identifier: ExtSystem,
				System:     "aa",
			},
			want2: "bb",
		},
		{
			name:    "error parsing EntityValue",
			source:  `'`,
			wantErr: true,
		},
		{
			name:    "error parsing ExternalID",
			source:  `SYSTEM`,
			wantErr: true,
		},
		{
			name:    "error empty",
			source:  ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, got1, got2, err := p.parseEntityDef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEntityDef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseEntityDef() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Parser.parseEntityDef() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Parser.parseEntityDef() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestParser_parsePEDef(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    EntityValue
		want1   *ExternalID
		wantErr bool
	}{
		{
			name:   "EntityValue",
			source: `'ab%aa;'`,
			want: EntityValue{
				"ab",
				&PERef{
					Name: "aa",
				},
			},
		},
		{
			name:   "ExternalID",
			source: `SYSTEM "aa"`,
			want1: &ExternalID{
				Identifier: ExtSystem,
				System:     "aa",
			},
		},
		{
			name:    "error parsing EntityValue",
			source:  `'`,
			wantErr: true,
		},
		{
			name:    "error parsing ExternalID",
			source:  `SYSTEM`,
			wantErr: true,
		},
		{
			name:    "error empty",
			source:  ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, got1, err := p.parsePEDef()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parsePEDef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parsePEDef() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Parser.parsePEDef() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParser_parseExternalID(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *ExternalID
		wantErr bool
	}{
		{
			name:    "not starts with SYSTEM or PUBLIC",
			source:  `SYS "aaa"`,
			wantErr: true,
		},
		{
			name:    "no space",
			source:  `SYSTEM"aaa"`,
			wantErr: true,
		},
		{
			name:    "pubid literal error",
			source:  `PUBLIC aa"`,
			wantErr: true,
		},
		{
			name:    "no space after pubid literal",
			source:  `PUBLIC "aa"`,
			wantErr: true,
		},
		{
			name:    "system literal error",
			source:  `SYSTEM aa"`,
			wantErr: true,
		},
		{
			name:   "public no error",
			source: `PUBLIC "pub" "sys"`,
			want: &ExternalID{
				Identifier: ExtPublic,
				Pubid:      "pub",
				System:     "sys",
			},
		},
		{
			name:   "system no error",
			source: `SYSTEM "sys"`,
			want: &ExternalID{
				Identifier: ExtSystem,
				System:     "sys",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseExternalID()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseExternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseNData(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "no space",
			source:  "NDATA aaa",
			wantErr: true,
		},
		{
			name:    "no NDATS",
			source:  "  aaa",
			wantErr: true,
		},
		{
			name:    "no space after NDATA",
			source:  " NDATAaaa",
			wantErr: true,
		},
		{
			source: " NDATA aaa",
			want:   "aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseNData()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseNData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseNData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEncoding(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "not starts with spaces",
			source:  `encoding="UTF-8"`,
			wantErr: true,
		},
		{
			name:    "typo encoding",
			source:  ` enco="UTF-8"`,
			wantErr: true,
		},
		{
			name:    "not equal",
			source:  ` encoding:"UTF-8"`,
			wantErr: true,
		},
		{
			name:    "no quote",
			source:  ` encoding=UTF-8`,
			wantErr: true,
		},
		{
			name:    "error while parsing encoding name",
			source:  ` encoding="あ" `,
			wantErr: true,
		},
		{
			name:    "different quote",
			source:  ` encoding="UTF-8' `,
			wantErr: true,
		},
		{
			source: ` encoding="UTF-8" `,
			want:   "UTF-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEncoding()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseEncName(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "not starts with alphabet",
			source:  "8UTF",
			wantErr: true,
		},
		{
			source: "UTF-8",
			want:   "UTF-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseEncName()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseEncName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.parseEncName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseNotation(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *Notation
		wantErr bool
	}{
		{
			name:    "not starts with <!NOTATION",
			source:  `<NOTATION name PUBLIC "public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "no space between <!NOTATION and name",
			source:  `<!NOTATIONname PUBLIC "public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "name error",
			source:  `<!NOTATION .name PUBLIC "public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "no space after name",
			source:  `<!NOTATION name'PUBLIC "public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "invalid external identifier",
			source:  `<!NOTATION name PUB "public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "no space external identifier",
			source:  `<!NOTATION name PUBLIC"public_ID" "URI">`,
			wantErr: true,
		},
		{
			name:    "error parsing pubid literal",
			source:  `<!NOTATION name PUBLIC " >`,
			wantErr: true,
		},
		{
			name:    "error parsing system literal",
			source:  `<!NOTATION name SYSTEM " >`,
			wantErr: true,
		},
		{
			name:    "not closed",
			source:  `<!NOTATION name PUBLIC "public_ID" "URI"`,
			wantErr: true,
		},
		{
			name:   "PublicID with system",
			source: `<!NOTATION name PUBLIC "public_ID" "URI">`,
			want: &Notation{
				Name: "name",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "public_ID",
					System:     "URI",
				},
			},
		},
		{
			name:   "PublicID without system",
			source: `<!NOTATION name PUBLIC "public_ID">`,
			want: &Notation{
				Name: "name",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "public_ID",
				},
			},
		},
		{
			name:   "PublicID without system, has space after pubid",
			source: `<!NOTATION name PUBLIC "public_ID" >`,
			want: &Notation{
				Name: "name",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "public_ID",
				},
			},
		},
		{
			name:   "System",
			source: `<!NOTATION name SYSTEM "system">`,
			want: &Notation{
				Name: "name",
				ExtID: ExternalID{
					Identifier: ExtSystem,
					System:     "system",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.source)
			got, err := p.parseNotation()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseNotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
