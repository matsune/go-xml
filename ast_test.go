package xml

import (
	"testing"
)

func TestXMLDecl_String(t *testing.T) {
	type fields struct {
		Version    string
		Encoding   string
		Standalone bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Version:    "1.0",
				Encoding:   "UTF-8",
				Standalone: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`,
		},
		{
			fields: fields{
				Version: "1.1",
			},
			want: `<?xml version="1.1" standalone="no" ?>`,
		},
		{
			fields: fields{},
			want:   `<?xml version="" standalone="no" ?>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XMLDecl{
				Version:    tt.fields.Version,
				Encoding:   tt.fields.Encoding,
				Standalone: tt.fields.Standalone,
			}
			if got := x.String(); got != tt.want {
				t.Errorf("XMLDecl.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDOCType_String(t *testing.T) {
	tests := []struct {
		name string
		DOCType
		want string
	}{
		{
			name: "only name",
			DOCType: DOCType{
				Name: "html",
			},
			want: `<!DOCTYPE html>`,
		},
		{
			name: "name and ExtID",
			DOCType: DOCType{
				Name: "html",
				ExtID: &ExternalID{
					Identifier: ExtPublic,
					Pubid:      "pubid",
					System:     "system",
				},
			},
			want: `<!DOCTYPE html PUBLIC "pubid" "system">`,
		},
		{
			name: "name and markups",
			DOCType: DOCType{
				Name: "html",
				Markups: []Markup{
					Comment("comment"),
				},
			},
			want: `<!DOCTYPE html [<!--comment-->]>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.DOCType.String(); got != tt.want {
				t.Errorf("DOCType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalID_String(t *testing.T) {
	type fields struct {
		Identifier ExtIdent
		Pubid      string
		System     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "system",
			fields: fields{
				Identifier: ExtSystem,
				System:     "system",
			},
			want: `SYSTEM "system"`,
		},
		{
			name: "pubid",
			fields: fields{
				Identifier: ExtPublic,
				Pubid:      "pubid",
			},
			want: `PUBLIC "pubid"`,
		},
		{
			name: "pubid and system",
			fields: fields{
				Identifier: ExtPublic,
				Pubid:      "pubid",
				System:     "system",
			},
			want: `PUBLIC "pubid" "system"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ExternalID{
				Identifier: tt.fields.Identifier,
				Pubid:      tt.fields.Pubid,
				System:     tt.fields.System,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("ExternalID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttlist_String(t *testing.T) {
	type fields struct {
		Name string
		Defs []*AttDef
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Name: "task",
				Defs: []*AttDef{
					&AttDef{
						Name: "status",
						Type: &Enum{
							Cases: []string{"important", "normal"},
						},
						Decl: &DefaultDecl{
							Type: IMPLIED,
						},
					},
				},
			},
			want: `<!ATTLIST task status (important|normal) #IMPLIED>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attlist{
				Name: tt.fields.Name,
				Defs: tt.fields.Defs,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("Attlist.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_String(t *testing.T) {
	type fields struct {
		Name  string
		Type  EntityType
		Value EntityValue
		ExtID *ExternalID
		NData string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "GEDecl, EntityValue",
			fields: fields{
				Name: "name",
				Type: EntityType_GE,
				Value: EntityValue{
					"value",
					PERef{
						Name: "peref",
					},
				},
			},
			want: `<!ENTITY name "value%peref;">`,
		},
		{
			name: "GEDecl, ExternalID, NData",
			fields: fields{
				Name: "name",
				Type: EntityType_GE,
				ExtID: &ExternalID{
					Identifier: ExtPublic,
					Pubid:      "pubid",
					System:     "system",
				},
				NData: "ndata",
			},
			want: `<!ENTITY name PUBLIC "pubid" "system" NDATA ndata>`,
		},
		{
			name: "PEDecl, ExternalID",
			fields: fields{
				Name: "name",
				Type: EntityType_PE,
				ExtID: &ExternalID{
					Identifier: ExtSystem,
					System:     "system",
				},
			},
			want: `<!ENTITY % name SYSTEM "system">`,
		},
		{
			name: "PEDecl, EntityValue",
			fields: fields{
				Name: "name",
				Type: EntityType_PE,
				Value: EntityValue{
					"string",
					PERef{
						Name: "peref",
					},
					CharRef{
						Prefix: "&#",
						Value:  "0",
					},
					EntityRef{
						Name: "entityref",
					},
				},
			},
			want: `<!ENTITY % name "string%peref;&#0;&entityref;">`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entity{
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
				ExtID: tt.fields.ExtID,
				NData: tt.fields.NData,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("Entity.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotation_String(t *testing.T) {
	type fields struct {
		Name  string
		ExtID ExternalID
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Name: "nota",
				ExtID: ExternalID{
					Identifier: ExtSystem,
					System:     "system",
				},
			},
			want: `<!NOTATION nota SYSTEM "system">`,
		},
		{
			fields: fields{
				Name: "nota",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "pubid",
				},
			},
			want: `<!NOTATION nota PUBLIC "pubid">`,
		},
		{
			fields: fields{
				Name: "nota",
				ExtID: ExternalID{
					Identifier: ExtPublic,
					Pubid:      "pubid",
					System:     "system",
				},
			},
			want: `<!NOTATION nota PUBLIC "pubid" "system">`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Notation{
				Name:  tt.fields.Name,
				ExtID: tt.fields.ExtID,
			}
			if got := n.String(); got != tt.want {
				t.Errorf("Notation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPI_String(t *testing.T) {
	type fields struct {
		Target      string
		Instruction string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Target:      "target",
				Instruction: "",
			},
			want: `<?target?>`,
		},
		{
			fields: fields{
				Target:      "target",
				Instruction: "inst",
			},
			want: `<?target inst?>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PI{
				Target:      tt.fields.Target,
				Instruction: tt.fields.Instruction,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("PI.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_String(t *testing.T) {
	tests := []struct {
		name string
		c    Comment
		want string
	}{
		{
			c:    Comment(" this is a comment "),
			want: `<!-- this is a comment -->`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Comment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttDef_String(t *testing.T) {
	type fields struct {
		Name string
		Type AttType
		Decl *DefaultDecl
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AttDef{
				Name: tt.fields.Name,
				Type: tt.fields.Type,
				Decl: tt.fields.Decl,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("AttDef.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDecl_String(t *testing.T) {
	type fields struct {
		Type     DefaultDeclType
		AttValue AttValue
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Type: REQUIRED,
			},
			want: "#REQUIRED",
		},
		{
			fields: fields{
				Type: IMPLIED,
			},
			want: "#IMPLIED",
		},
		{
			fields: fields{
				Type: FIXED,
				AttValue: AttValue{
					"a", &EntityRef{
						Name: "entity",
					},
				},
			},
			want: `#FIXED "a&entity;"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultDecl{
				Type:     tt.fields.Type,
				AttValue: tt.fields.AttValue,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("DefaultDecl.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotationType_String(t *testing.T) {
	type fields struct {
		Names []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Names: []string{
					"a", "b", "c",
				},
			},
			want: `NOTATION (a|b|c)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NotationType{
				Names: tt.fields.Names,
			}
			if got := n.String(); got != tt.want {
				t.Errorf("NotationType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnum_String(t *testing.T) {
	type fields struct {
		Cases []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Cases: []string{
					"a", "b", "c",
				},
			},
			want: `(a|b|c)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Enum{
				Cases: tt.fields.Cases,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("Enum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChildren_String(t *testing.T) {
	type fields struct {
		ChoiceSeq ChoiceSeq
		Suffix    *rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				ChoiceSeq: &Seq{
					CPs: []CP{
						CP{
							Name: "a",
						},
						CP{
							Name: "b",
						},
						CP{
							ChoiceSeq: &Choice{
								CPs: []CP{
									CP{
										Name: "c",
									},
									CP{
										Name: "d",
									},
								},
							},
						},
					},
				},
				Suffix: newRune('?'),
			},
			want: `(a,b,(c|d))?`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Children{
				ChoiceSeq: tt.fields.ChoiceSeq,
				Suffix:    tt.fields.Suffix,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Children.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCP_String(t *testing.T) {
	type fields struct {
		Name      string
		ChoiceSeq ChoiceSeq
		Suffix    *rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "name",
			fields: fields{
				Name:   "n",
				Suffix: newRune('?'),
			},
			want: `n?`,
		},
		{
			name: "choice",
			fields: fields{
				ChoiceSeq: &Choice{
					CPs: []CP{
						CP{
							Name: "a",
						},
						CP{
							Name: "b",
						},
					},
				},
				Suffix: newRune('*'),
			},
			want: `(a|b)*`,
		},
		{
			name: "seq",
			fields: fields{
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
				Suffix: newRune('+'),
			},
			want: `(a,b)+`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CP{
				Name:      tt.fields.Name,
				ChoiceSeq: tt.fields.ChoiceSeq,
				Suffix:    tt.fields.Suffix,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("CP.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChoice_String(t *testing.T) {
	type fields struct {
		CPs []CP
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				CPs: []CP{
					CP{
						Name: "a",
					},
					CP{
						Name: "b",
					},
				},
			},
			want: "(a|b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Choice{
				CPs: tt.fields.CPs,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Choice.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_String(t *testing.T) {
	type fields struct {
		CPs []CP
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				CPs: []CP{
					CP{
						Name: "a",
					},
					CP{
						Name: "b",
					},
				},
			},
			want: "(a,b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Seq{
				CPs: tt.fields.CPs,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("Seq.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityValue_String(t *testing.T) {
	tests := []struct {
		name string
		e    EntityValue
		want string
	}{
		{
			e: EntityValue{
				"string",
				PERef{
					Name: "peref",
				},
				CharRef{
					Prefix: "&#",
					Value:  "0",
				},
				EntityRef{
					Name: "entity",
				},
			},
			want: `"string%peref;&#0;&entity;"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EntityValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttValue_String(t *testing.T) {
	tests := []struct {
		name string
		a    AttValue
		want string
	}{
		{
			a: AttValue{
				"string",
				CharRef{
					Prefix: "&#",
					Value:  "0",
				},
				EntityRef{
					Name: "entity",
				},
			},
			want: `"string&#0;&entity;"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("AttValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharRef_String(t *testing.T) {
	type fields struct {
		Prefix string
		Value  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Prefix: "&#",
				Value:  "00",
			},
			want: `&#00;`,
		},
		{
			fields: fields{
				Prefix: "&#x",
				Value:  "0aF",
			},
			want: `&#x0aF;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CharRef{
				Prefix: tt.fields.Prefix,
				Value:  tt.fields.Value,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("CharRef.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityRef_String(t *testing.T) {
	tests := []struct {
		name    string
		refName string
		want    string
	}{
		{
			refName: "aa",
			want:    `&aa;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EntityRef{
				Name: tt.refName,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("EntityRef.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPERef_String(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name    string
		refName string
		want    string
	}{
		{
			refName: "peref",
			want:    `%peref;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := PERef{
				Name: tt.refName,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("PERef.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCData_String(t *testing.T) {
	tests := []struct {
		name string
		c    CData
		want string
	}{
		{
			c:    CData("cdata"),
			want: `<![CDATA[cdata]]>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("CData.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
