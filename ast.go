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
		Markups     []Markup
		PEReference string
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

type (
	Markup interface {
		Markup()
	}

	Element struct {
		Name string
		ContentSpec
	}
	Attlist  struct{}
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

// ContentSpec
type (
	ContentSpec interface {
		ContentSpec()
	}

	EMPTY struct{}
	ANY   struct{}
	Mixed struct {
		Names []string
	}
	Choice struct {
		Names []string // separated '|'
	}
	Seq struct {
		Names []string // separated ','
	}
)

func (EMPTY) ContentSpec()  {}
func (ANY) ContentSpec()    {}
func (Mixed) ContentSpec()  {}
func (Choice) ContentSpec() {}
func (Seq) ContentSpec()    {}
