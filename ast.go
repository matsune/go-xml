package xml

import (
	"fmt"
)

type (
	Prolog struct {
		*XMLDecl
		*DOCType
	}

	XMLDecl struct {
		Version    string
		Encoding   string
		Standalone bool
	}

	DOCType struct {
		Name string
		*ExternalID
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

const (
	ExtSystem ExtIdent = "SYSTEM"
	ExtPublic          = "PUBLIC"
)

// Markup
type (
	Markup interface {
		Markup()
	}

	Element struct {
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
		ExID  *ExternalID
		NData string
	}
	Notation struct{}
	PI       struct{}
	Comment  string
)

const (
	EntityType_GE EntityType = iota
	EntityType_PE
)

func (Element) Markup()  {}
func (Attlist) Markup()  {}
func (Entity) Markup()   {}
func (Notation) Markup() {}
func (PI) Markup()       {}
func (Comment) Markup()  {}

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
