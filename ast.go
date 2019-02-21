package xml

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
		PERef   PERef
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
		Defs []AttDef
	}
	Entity   struct{}
	Notation struct{}
	PI       struct{}
	Comment  string
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
		Decl DefaultDecl
	}

	DefaultDeclType string
	DefaultDecl     struct {
		Type DefaultDeclType
		Refs []Ref
	}

	NotationType struct {
		Types []string
	}

	Enum struct {
		Nms []string
	}
)

const (
	REQUIRED DefaultDeclType = "#REQUIRED"
	IMPLIED                  = "#IMPLIED"
	FIXED                    = "#FIXED"
)

const (
	Att_CDATA    StringType    = "CDATA"
	Att_ID       TokenizedType = "ID"
	Att_IDREF                  = "IDREF"
	Att_IDREFS                 = "IDREFS"
	Att_ENTITY                 = "ENTITY"
	Att_ENTITIES               = "ENTITIES"
	Att_NMTOKEN                = "NMTOKEN"
	Att_NMTOKENS               = "NMTOKENS"
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

// Entity Ref
type (
	Ref interface {
		Ref()
	}
	EntityRef string // & Name ;
	PERef     string // % Name ;
)

func (EntityRef) Ref() {}
func (PERef) Ref()     {}
