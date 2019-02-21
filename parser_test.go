package xml

import (
	"fmt"
	"reflect"
	"testing"
)

func newRune(r rune) *rune { return &r }

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

func TestParser_parseDoctype(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    *DOCType
		wantErr bool
	}{
		{
			name:    "no '!'",
			source:  `<DOCTYPE document SYSTEM "subjects.dtd">`,
			wantErr: true,
		},
		{
			name:    "no space between DOCTYPE and name",
			source:  `<!DOCTYPEdocument SYSTEM "subjects.dtd">`,
			wantErr: true,
		},
		{
			name:    "no name",
			source:  `<!DOCTYPE >`,
			wantErr: true,
		},
		{
			name:    "no system literal",
			source:  `<!DOCTYPE doc SYSTEM>`,
			wantErr: true,
		},
		{
			name:    "no end tag",
			source:  `<!DOCTYPE doc SYSTEM "subjects.dtd" `,
			wantErr: true,
		},
		{
			name:   "has ExternalID",
			source: `<!DOCTYPE doc SYSTEM "subjects.dtd" >`,
			want: &DOCType{
				Name: "doc",
				ExternalID: &ExternalID{
					Identifier: ExtSystem,
					System:     "subjects.dtd",
				},
			},
		},
		{
			source: `<!DOCTYPE student [
				<!ELEMENT student (surname,firstname*,dob?,(origin|sex)?)>
				<!ELEMENT surname (#PCDATA)>
				<!ELEMENT firstname (#PCDATA)>
				<!ELEMENT sex (#PCDATA)>
			]>`,
			want: &DOCType{
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
			source:  ` stand=yes`,
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

func TestParser_parseEncoding(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:    "not starts with spaces",
			source:  `encoding = "UTF-8"`,
			wantErr: true,
		},
		{
			name:    "typo encoding",
			source:  ` encod = "UTF-8"`,
			wantErr: true,
		},
		{
			name:    "not equal",
			source:  ` encoding:"UTF-8"`,
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
