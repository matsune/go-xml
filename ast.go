package xml

import (
	"fmt"
)

type (
	XML struct {
		*Prolog
		*Element
		Misc []interface{}
	}
)

type (
	Prolog struct {
		// ignoring Miscs
		*XMLDecl
		*DOCType
	}

	XMLDecl struct {
		Version    string
		Encoding   string
		Standalone bool
	}

	DOCType struct {
		Name    string
		ExtID   *ExternalID
		Markups []Markup
		PERef   *PERef
	}

	ExtIdent string

	ExternalID struct {
		Identifier ExtIdent
		Pubid      string
		System     string
	}
)

func (x XMLDecl) String() string {
	str := "<?xml"

	str += fmt.Sprintf(` version="%s"`, x.Version)

	if len(x.Encoding) > 0 {
		str += fmt.Sprintf(` encoding="%s"`, x.Encoding)
	}

	stdStr := "no"
	if x.Standalone {
		stdStr = "yes"
	}
	str += fmt.Sprintf(` standalone="%s" ?>`, stdStr)

	return str
}

func (d DOCType) String() string {
	str := "<!DOCTYPE"

	str += fmt.Sprintf(` %s`, d.Name)

	if d.ExtID != nil {
		str += fmt.Sprintf(" %v", d.ExtID)
	}

	if len(d.Markups) > 0 {
		str += " ["
	}
	for i, m := range d.Markups {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf("%v", m)
	}
	if len(d.Markups) > 0 {
		str += "]"
	}

	str += ">"
	return str
}

const (
	ExtSystem ExtIdent = "SYSTEM"
	ExtPublic          = "PUBLIC"
)

func (e ExternalID) String() string {
	str := string(e.Identifier)
	if len(e.Pubid) > 0 {
		str += fmt.Sprintf(" %q", e.Pubid)
	}
	if len(e.System) > 0 {
		str += fmt.Sprintf(" %q", e.System)
	}
	return str
}

// Markup
type (
	Markup interface {
		Markup()
	}

	ElementDecl struct {
		Name string
		ContentSpec
	}
	Attlist struct {
		Name string
		Defs []*AttDef
	}
	EntityType int
	Entity     struct {
		Name  string
		Type  EntityType
		Value EntityValue
		ExtID *ExternalID
		NData string
	}
	Notation struct {
		Name  string
		ExtID ExternalID
	}
	PI struct {
		Target      string
		Instruction string
	}
	Comment string
)

const (
	EntityType_GE EntityType = iota
	EntityType_PE
)

func (ElementDecl) Markup() {}
func (Attlist) Markup()     {}
func (Entity) Markup()      {}
func (Notation) Markup()    {}
func (PI) Markup()          {}
func (Comment) Markup()     {}

func (c Comment) String() string {
	return fmt.Sprintf("<!--%s-->", string(c))
}

// Attribute Types
type (
	AttType interface {
		AttType()
	}

	StringType    string
	TokenizedType string

	AttDef struct {
		Name string
		Type AttType
		Decl *DefaultDecl
	}

	DefaultDeclType string
	DefaultDecl     struct {
		Type     DefaultDeclType
		AttValue AttValue
	}

	NotationType struct {
		Names []string
	}

	Enum struct {
		Cases []string
	}
)

const (
	REQUIRED DefaultDeclType = "#REQUIRED"
	IMPLIED  DefaultDeclType = "#IMPLIED"
	FIXED    DefaultDeclType = "#FIXED"
)

const (
	Att_CDATA StringType = "CDATA"

	Att_ID       TokenizedType = "ID"
	Att_IDREF    TokenizedType = "IDREF"
	Att_IDREFS   TokenizedType = "IDREFS"
	Att_ENTITY   TokenizedType = "ENTITY"
	Att_ENTITIES TokenizedType = "ENTITIES"
	Att_NMTOKEN  TokenizedType = "NMTOKEN"
	Att_NMTOKENS TokenizedType = "NMTOKENS"
)

func (StringType) AttType()    {}
func (TokenizedType) AttType() {}
func (NotationType) AttType()  {}
func (Enum) AttType()          {}

// ContentSpec, ChoiseSeq
type (
	ContentSpec interface {
		ContentSpec()
	}

	EMPTY struct{}
	ANY   struct{}
	Mixed struct {
		Names []string
	}

	Children struct {
		ChoiceSeq
		Suffix *rune // null or '?' or '*' or '+'
	}
	CP struct {
		Name string
		ChoiceSeq
		Suffix *rune
	}

	ChoiceSeq interface {
		ChoiceSeq()
	}
	Choice struct {
		CPs []CP // separated '|'
	}
	Seq struct {
		CPs []CP // separated ','
	}
)

func (EMPTY) ContentSpec()    {}
func (ANY) ContentSpec()      {}
func (Mixed) ContentSpec()    {}
func (Children) ContentSpec() {}

func (Choice) ChoiceSeq() {}
func (Seq) ChoiceSeq()    {}

// Ref
type (
	Ref interface {
		Ref()
	}

	// string or PERef or CharRef or EntityRef
	EntityValue []interface{}
	// string or CharRef or EntityRef
	AttValue []interface{}

	CharRef struct {
		Prefix string // &# or &#x
		Value  string
	}
	EntityRef struct {
		Name string // & Name ;
	}
	PERef struct {
		Name string // % Name ;
	}
)

func (CharRef) Ref()   {}
func (EntityRef) Ref() {}

func (e CharRef) String() string {
	return fmt.Sprintf("%s%s;", e.Prefix, e.Value)
}
func (e EntityRef) String() string {
	return fmt.Sprintf("&%s;", e.Name)
}
func (e PERef) String() string {
	return "%" + fmt.Sprintf(`%s;`, e.Name)
}

// Element
type (
	// <name />
	Element struct {
		Name       string
		Attrs      Attributes
		Contents   []interface{}
		IsEmptyTag bool
	}

	Attribute struct {
		Name string
		AttValue
	}

	Attributes []*Attribute
)

type (
	CData string
)

func (c CData) String() string {
	return fmt.Sprintf("<![CDATA[%s]]>", string(c))
}
