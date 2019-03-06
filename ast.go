package xml

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
)

type (
	ExternalType int

	ExternalID struct {
		Type   ExternalType
		Pubid  string
		System string
	}
)

const (
	EXT_SYSTEM ExternalType = iota
	EXT_PUBLIC
)

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
	ENTITY_GE EntityType = iota
	ENTITY_PE
)

func (ElementDecl) Markup() {}
func (Attlist) Markup()     {}
func (Entity) Markup()      {}
func (Notation) Markup()    {}
func (PI) Markup()          {}
func (Comment) Markup()     {}

// Attribute Types
type (
	AttDef struct {
		Name string
		Type AttType
		Decl *DefaultDecl
	}

	AttType interface {
		AttType()
	}

	AttToken int

	NotationType struct {
		Names []string
	}

	Enum struct {
		Cases []string
	}

	DefaultDeclType int

	DefaultDecl struct {
		Type DefaultDeclType
		AttValue
	}
)

func (AttToken) AttType()     {}
func (NotationType) AttType() {}
func (Enum) AttType()         {}

const (
	_ AttToken = iota
	ATT_CDATA
	ATT_ID
	ATT_IDREF
	ATT_IDREFS
	ATT_ENTITY
	ATT_ENTITIES
	ATT_NMTOKEN
	ATT_NMTOKENS
)

const (
	_ DefaultDeclType = iota
	DECL_REQUIRED
	DECL_IMPLIED
	DECL_FIXED
)

// ContentSpec, ChoiseSeq
type (
	ContentSpec interface {
		ContentSpec()
	}

	EMPTY struct{} // EMPTY
	ANY   struct{} // ANY
	Mixed struct { // #PCDATA
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

// Element
type (
	Element struct {
		Name       string
		Attrs      Attributes
		Contents   []interface{} // PERef, CDSect, Comment, PI, Element, String
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
