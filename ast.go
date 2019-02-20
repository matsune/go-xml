package xml

type (
	Node interface {
		Node()
	}
)

type (
	XMLDecl struct {
		VersionInfo  string
		EncodingDecl string
		Standalone   bool
	}
)

func (XMLDecl) Node() {}
