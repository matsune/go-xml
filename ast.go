package xml

type (
	AST interface {
		AST()
	}

	Terminal interface {
		AST
		ToString() string
	}
)

type (
	XML struct {
		*Prolog
		*Element
		Misc []interface{}
	}

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

	ExternalType int

	ExternalID struct {
		Type   ExternalType
		Pubid  string
		System string
	}

	// Markup >>

	Markup interface {
		Terminal
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

	Entity struct {
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

	// Attribute Types

	AttDef struct {
		Name string
		Type AttType
		Decl *DefaultDecl
	}

	AttType interface {
		Terminal
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

	// ContentSpec, ChoiceSeq

	ContentSpec interface {
		Terminal
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
		Terminal
		ChoiceSeq()
	}
	Choice struct {
		CPs []CP // separated '|'
	}
	Seq struct {
		CPs []CP // separated ','
	}

	// Ref

	Ref interface {
		Terminal
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

	// Element

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

	CData string
)

// XML, Prolog, XMLDecl, DOCType and Element are Non-Terminal
func (XML) AST()             {}
func (Prolog) AST()          {}
func (XMLDecl) AST()         {}
func (DOCType) AST()         {}
func (ExternalType) AST()    {}
func (ExternalID) AST()      {}
func (ElementDecl) AST()     {}
func (Attlist) AST()         {}
func (EntityType) AST()      {}
func (Entity) AST()          {}
func (Notation) AST()        {}
func (PI) AST()              {}
func (Comment) AST()         {}
func (AttDef) AST()          {}
func (AttToken) AST()        {}
func (NotationType) AST()    {}
func (Enum) AST()            {}
func (DefaultDeclType) AST() {}
func (DefaultDecl) AST()     {}
func (EMPTY) AST()           {}
func (ANY) AST()             {}
func (Mixed) AST()           {}
func (Children) AST()        {}
func (CP) AST()              {}
func (Choice) AST()          {}
func (Seq) AST()             {}
func (EntityValue) AST()     {}
func (AttValue) AST()        {}
func (CharRef) AST()         {}
func (EntityRef) AST()       {}
func (PERef) AST()           {}
func (Element) AST()         {}
func (Attribute) AST()       {}
func (Attributes) AST()      {}
func (CData) AST()           {}

func (ElementDecl) Markup() {}
func (Attlist) Markup()     {}
func (Entity) Markup()      {}
func (Notation) Markup()    {}
func (PI) Markup()          {}
func (Comment) Markup()     {}

func (AttToken) AttType()     {}
func (NotationType) AttType() {}
func (Enum) AttType()         {}

func (EMPTY) ContentSpec()    {}
func (ANY) ContentSpec()      {}
func (Mixed) ContentSpec()    {}
func (Children) ContentSpec() {}

func (Choice) ChoiceSeq() {}
func (Seq) ChoiceSeq()    {}

func (CharRef) Ref()   {}
func (EntityRef) Ref() {}

const (
	ExternalTypeSystem ExternalType = iota
	ExternalTypePublic
)
const (
	EntityTypeGE EntityType = iota
	EntityTypePE
)
const (
	_ AttToken = iota
	AttTokenCDATA
	AttTokenID
	AttTokenIDREF
	AttTokenIDREFS
	AttTokenENTITY
	AttTokenENTITIES
	AttTokenNMTOKEN
	AttTokenNMTOKENS
)
const (
	_ DefaultDeclType = iota
	DefaultDeclTypeRequired
	DefaultDeclTypeImplied
	DefaultDeclTypeFixed
)
