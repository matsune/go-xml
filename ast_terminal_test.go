package xml

import (
	"testing"
)

func TestExternalType_String(t *testing.T) {
	tests := []struct {
		name string
		e    ExternalType
		want string
	}{
		{
			name: "public",
			e:    ExternalTypePublic,
			want: "PUBLIC",
		},
		{
			name: "system",
			e:    ExternalTypeSystem,
			want: "SYSTEM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ToString(); got != tt.want {
				t.Errorf("ExternalType.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalID_String(t *testing.T) {
	type fields struct {
		Type   ExternalType
		Pubid  string
		System string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "system",
			fields: fields{
				Type:   ExternalTypeSystem,
				System: "system",
			},
			want: `SYSTEM "system"`,
		},
		{
			name: "pubid",
			fields: fields{
				Type:  ExternalTypePublic,
				Pubid: "pubid",
			},
			want: `PUBLIC "pubid"`,
		},
		{
			name: "pubid and system",
			fields: fields{
				Type:   ExternalTypePublic,
				Pubid:  "pubid",
				System: "system",
			},
			want: `PUBLIC "pubid" "system"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ExternalID{
				Type:   tt.fields.Type,
				Pubid:  tt.fields.Pubid,
				System: tt.fields.System,
			}
			if got := e.ToString(); got != tt.want {
				t.Errorf("ExternalID.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElementDecl_String(t *testing.T) {
	type fields struct {
		Name        string
		ContentSpec ContentSpec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "contentSpce is EMPTY",
			fields: fields{
				Name:        "name",
				ContentSpec: &EMPTY{},
			},
			want: `<!ELEMENT name EMPTY>`,
		},
		{
			name: "contentSpce is ANY",
			fields: fields{
				Name:        "name",
				ContentSpec: &ANY{},
			},
			want: `<!ELEMENT name ANY>`,
		},
		{
			name: "contentSpce is Mixed",
			fields: fields{
				Name: "name",
				ContentSpec: &Mixed{
					Names: []string{"a", "b"},
				},
			},
			want: `<!ELEMENT name (#PCDATA|a|b)>`,
		},
		{
			name: "contentSpce is Children",
			fields: fields{
				Name: "name",
				ContentSpec: &Children{
					ChoiceSeq: &Seq{
						CPs: []CP{
							CP{
								Name: "a",
							},
						},
					},
					Suffix: newRune('*'),
				},
			},
			want: `<!ELEMENT name (a)*>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ElementDecl{
				Name:        tt.fields.Name,
				ContentSpec: tt.fields.ContentSpec,
			}
			if got := e.ToString(); got != tt.want {
				t.Errorf("ElementDecl.ToString() = %v, want %v", got, tt.want)
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
							Type: DefaultDeclTypeImplied,
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
			if got := a.ToString(); got != tt.want {
				t.Errorf("Attlist.ToString() = %v, want %v", got, tt.want)
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
				Type: EntityTypeGE,
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
				Type: EntityTypeGE,
				ExtID: &ExternalID{
					Type:   ExternalTypePublic,
					Pubid:  "pubid",
					System: "system",
				},
				NData: "ndata",
			},
			want: `<!ENTITY name PUBLIC "pubid" "system" NDATA ndata>`,
		},
		{
			name: "PEDecl, ExternalID",
			fields: fields{
				Name: "name",
				Type: EntityTypePE,
				ExtID: &ExternalID{
					Type:   ExternalTypeSystem,
					System: "system",
				},
			},
			want: `<!ENTITY % name SYSTEM "system">`,
		},
		{
			name: "PEDecl, EntityValue",
			fields: fields{
				Name: "name",
				Type: EntityTypePE,
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
			if got := e.ToString(); got != tt.want {
				t.Errorf("Entity.ToString() = %v, want %v", got, tt.want)
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
					Type:   ExternalTypeSystem,
					System: "system",
				},
			},
			want: `<!NOTATION nota SYSTEM "system">`,
		},
		{
			fields: fields{
				Name: "nota",
				ExtID: ExternalID{
					Type:  ExternalTypePublic,
					Pubid: "pubid",
				},
			},
			want: `<!NOTATION nota PUBLIC "pubid">`,
		},
		{
			fields: fields{
				Name: "nota",
				ExtID: ExternalID{
					Type:   ExternalTypePublic,
					Pubid:  "pubid",
					System: "system",
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
			if got := n.ToString(); got != tt.want {
				t.Errorf("Notation.ToString() = %v, want %v", got, tt.want)
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
			if got := p.ToString(); got != tt.want {
				t.Errorf("PI.ToString() = %v, want %v", got, tt.want)
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
			if got := tt.c.ToString(); got != tt.want {
				t.Errorf("Comment.ToString() = %v, want %v", got, tt.want)
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
			if got := a.ToString(); got != tt.want {
				t.Errorf("AttDef.ToString() = %v, want %v", got, tt.want)
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
				Type: DefaultDeclTypeRequired,
			},
			want: "#REQUIRED",
		},
		{
			fields: fields{
				Type: DefaultDeclTypeImplied,
			},
			want: "#IMPLIED",
		},
		{
			fields: fields{
				Type: DefaultDeclTypeFixed,
				AttValue: AttValue{
					"a",
					&EntityRef{
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
			if got := d.ToString(); got != tt.want {
				t.Errorf("DefaultDecl.ToString() = %v, want %v", got, tt.want)
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
			if got := n.ToString(); got != tt.want {
				t.Errorf("NotationType.ToString() = %v, want %v", got, tt.want)
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
			if got := e.ToString(); got != tt.want {
				t.Errorf("Enum.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDeclType_String(t *testing.T) {
	tests := []struct {
		name string
		d    DefaultDeclType
		want string
	}{
		{
			d:    DefaultDeclTypeRequired,
			want: "#REQUIRED",
		},
		{
			d:    DefaultDeclTypeImplied,
			want: "#IMPLIED",
		},
		{
			d:    DefaultDeclTypeFixed,
			want: "#FIXED",
		},
		{
			name: "unknown",
			d:    0,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ToString(); got != tt.want {
				t.Errorf("DefaultDeclType.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttToken_String(t *testing.T) {
	tests := []struct {
		name string
		a    AttToken
		want string
	}{
		{
			a:    AttTokenCDATA,
			want: "CDATA",
		},
		{
			a:    AttTokenID,
			want: "ID",
		},
		{
			a:    AttTokenIDREF,
			want: "IDREF",
		},
		{
			a:    AttTokenIDREFS,
			want: "IDREFS",
		},
		{
			a:    AttTokenENTITY,
			want: "ENTITY",
		},
		{
			a:    AttTokenENTITIES,
			want: "ENTITIES",
		},
		{
			a:    AttTokenNMTOKEN,
			want: "NMTOKEN",
		},
		{
			a:    AttTokenNMTOKENS,
			want: "NMTOKENS",
		},
		{
			name: "unknown",
			a:    0,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.ToString(); got != tt.want {
				t.Errorf("AttToken.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEMPTY_String(t *testing.T) {
	t.Run("EMPTY ToString()", func(t *testing.T) {
		want := "EMPTY"
		e := EMPTY{}
		if got := e.ToString(); got != want {
			t.Errorf("EMPTY.ToString() = %v, want %v", got, want)
		}
	})
}

func TestANY_String(t *testing.T) {
	t.Run("ANY ToString()", func(t *testing.T) {
		want := "ANY"
		e := ANY{}
		if got := e.ToString(); got != want {
			t.Errorf("EMPTY.ToString() = %v, want %v", got, want)
		}
	})
}

func TestMixed_String(t *testing.T) {
	type fields struct {
		Names []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty names",
			want: `(#PCDATA)`,
		},
		{
			fields: fields{
				Names: []string{"a", "b"},
			},
			want: `(#PCDATA|a|b)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Mixed{
				Names: tt.fields.Names,
			}
			if got := m.ToString(); got != tt.want {
				t.Errorf("Mixed.ToString() = %v, want %v", got, tt.want)
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
			if got := c.ToString(); got != tt.want {
				t.Errorf("Children.ToString() = %v, want %v", got, tt.want)
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
			if got := c.ToString(); got != tt.want {
				t.Errorf("CP.ToString() = %v, want %v", got, tt.want)
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
			if got := c.ToString(); got != tt.want {
				t.Errorf("Choice.ToString() = %v, want %v", got, tt.want)
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
			if got := s.ToString(); got != tt.want {
				t.Errorf("Seq.ToString() = %v, want %v", got, tt.want)
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
			if got := tt.e.ToString(); got != tt.want {
				t.Errorf("EntityValue.ToString() = %v, want %v", got, tt.want)
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
			if got := tt.a.ToString(); got != tt.want {
				t.Errorf("AttValue.ToString() = %v, want %v", got, tt.want)
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
			if got := e.ToString(); got != tt.want {
				t.Errorf("CharRef.ToString() = %v, want %v", got, tt.want)
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
			if got := e.ToString(); got != tt.want {
				t.Errorf("EntityRef.ToString() = %v, want %v", got, tt.want)
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
			if got := e.ToString(); got != tt.want {
				t.Errorf("PERef.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttribute_String(t *testing.T) {
	type fields struct {
		Name     string
		AttValue AttValue
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Name: "attr",
				AttValue: AttValue{
					"string",
					&CharRef{
						Prefix: "&#",
						Value:  "0",
					},
					&EntityRef{
						Name: "entity",
					},
				},
			},
			want: `attr="string&#0;&entity;"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attribute{
				Name:     tt.fields.Name,
				AttValue: tt.fields.AttValue,
			}
			if got := a.ToString(); got != tt.want {
				t.Errorf("Attribute.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_String(t *testing.T) {
	tests := []struct {
		name string
		a    Attributes
		want string
	}{
		{
			a: Attributes{
				&Attribute{
					Name: "a",
					AttValue: AttValue{
						"str1",
					},
				},
				&Attribute{
					Name: "b",
					AttValue: AttValue{
						"str2",
					},
				},
				&Attribute{
					Name: "c",
					AttValue: AttValue{
						"str3",
						&CharRef{
							Prefix: "&#",
							Value:  "0",
						},
					},
				},
			},
			want: `a="str1" b="str2" c="str3&#0;"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.ToString(); got != tt.want {
				t.Errorf("Attributes.ToString() = %v, want %v", got, tt.want)
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
			if got := tt.c.ToString(); got != tt.want {
				t.Errorf("CData.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
