package xml

import "fmt"

type Parser struct {
	source []rune
	cursor int
	stack  []int
}

func NewParser(str string) *Parser {
	return &Parser{
		source: []rune(str),
		cursor: 0,
	}
}

// document	::=	prolog element Misc*
func (p *Parser) Parse() {
	err := p.parseProlog()
	if err != nil {
		panic(err)
	}
	// TODO: parse element
}

// save cursor
func (p *Parser) push() {
	p.stack = append(p.stack, p.cursor)
}

// load last saved cursor
func (p *Parser) pop() {
	if len(p.stack) > 0 {
		p.cursor = p.stack[len(p.stack)-1]
	}
}

func (p *Parser) Test(r rune) bool {
	return p.Get() == r
}

func (p *Parser) Must(r rune) error {
	if !p.Test(r) {
		return fmt.Errorf("expected %q", r)
	}
	p.Step()
	return nil
}

// Check whether string on the cursor starts with str.
// This method doesn't proceed cursor, only check.
func (p *Parser) Tests(str string) bool {
	i := p.cursor
	e := i + len([]rune(str))
	if len(p.source) < e {
		return false
	}
	s := p.source[i:e]
	return string(s) == str
}

// Returns error if string on cursor doesn't match str.
// This method proceeds cursor if matching.
func (p *Parser) Musts(str string) error {
	if !p.Tests(str) {
		return fmt.Errorf("expected %q", str)
	}
	p.StepN(len(str))
	return nil
}

func (p *Parser) Step() {
	p.cursor++
}

func (p *Parser) StepN(n int) {
	p.cursor += n
}

const (
	EOF = 0
)

func (p *Parser) Get() rune {
	if p.isEnd() {
		return EOF
	}
	return rune(p.source[p.cursor])
}

func (p *Parser) isEnd() bool {
	return len(p.source) <= p.cursor
}

/// EBNF for XML 1.0
/// http://www.jelks.nu/XML/xmlebnf.html#NT-VersionInfo

func (p *Parser) parseProlog() error {
	if p.Tests("<?xml") {
		xmlDecl, err := p.parseXmlDecl()
		if err != nil {
			return err
		}
		fmt.Println(xmlDecl)
	}
	return nil
}

// XMLDecl ::= '<?xml' VersionInfo EncodingDecl? SDDecl? S? '?>'
func (p *Parser) parseXmlDecl() (*XMLDecl, error) {
	if err := p.Musts("<?xml"); err != nil {
		return nil, err
	}
	x := XMLDecl{}

	ver, err := p.parseVersion()
	if err != nil {
		return nil, err
	}
	x.VersionInfo = ver

	p.push()
	p.skipSpace()
	if p.Tests("encoding") {
		p.pop()

		enc, err := p.parseEncoding()
		if err != nil {
			return nil, err
		}
		x.EncodingDecl = enc
	} else {
		p.pop()
	}

	p.push()
	p.skipSpace()
	if p.Tests("standalone") {
		p.pop()

		std, err := p.parseStandalone()
		if err != nil {
			return nil, err
		}
		x.Standalone = std
	} else {
		p.pop()
	}

	p.skipSpace()
	if err := p.Musts("?>"); err != nil {
		return nil, err
	}

	return &x, nil
}

// VersionInfo ::= S 'version' Eq (' VersionNum ' | " VersionNum ")
func (p *Parser) parseVersion() (ver string, err error) {
	p.parseSpace()

	if err = p.Musts("version"); err != nil {
		return
	}
	if err = p.parseEq(); err != nil {
		return
	}

	err = p.parseQuote()
	if err != nil {
		return
	}

	ver, err = p.parseVersionNum()
	if err != nil {
		return
	}

	err = p.parseQuote()
	if err != nil {
		return
	}

	return
}

func (p *Parser) parseQuote() error {
	var err error
	if isQuote(p.Get()) {
		p.Step()
	} else {
		err = fmt.Errorf("expected ' or \"")
	}
	return err
}

// EncodingDecl ::= S 'encoding' Eq ('"' EncName  '"' |  "'" EncName "'" )
func (p *Parser) parseEncoding() (string, error) {
	if err := p.parseSpace(); err != nil {
		return "", err
	}
	if err := p.Musts("encoding"); err != nil {
		return "", err
	}
	if err := p.parseEq(); err != nil {
		return "", err
	}

	if err := p.parseQuote(); err != nil {
		return "", err
	}

	enc, err := p.parseEncName()
	if err != nil {
		return "", err
	}

	if err = p.parseQuote(); err != nil {
		return "", err
	}

	return enc, nil
}

// EncName ::= [A-Za-z] ([A-Za-z0-9._] | '-')*
func (p *Parser) parseEncName() (string, error) {
	var str string
	r := p.Get()
	if !isAlpha(r) {
		return "", fmt.Errorf("error while parsing encoding name")
	}
	str += string(r)
	p.Step()

	for {
		if isAlpha(p.Get()) || isNum(p.Get()) || p.Test('.') || p.Test('_') || p.Test('-') {
			str += string(p.Get())
			p.Step()
		} else {
			break
		}
	}

	return str, nil
}

// SDDecl ::= S 'standalone' Eq (("'" ('yes' | 'no') "'") | ('"' ('yes' | 'no') '"'))
func (p *Parser) parseStandalone() (bool, error) {
	if err := p.parseSpace(); err != nil {
		return false, err
	}
	if err := p.Musts("standalone"); err != nil {
		return false, err
	}
	if err := p.parseEq(); err != nil {
		return false, err
	}
	if err := p.parseQuote(); err != nil {
		return false, err
	}
	var std bool
	if p.Tests("yes") {
		std = true
		p.StepN(3)
	} else if p.Tests("no") {
		p.StepN(2)
	} else {
		return false, fmt.Errorf("error while parsing standalone")
	}
	if err := p.parseQuote(); err != nil {
		return false, err
	}
	return std, nil
}

// S ::= (#x20 | #x9 | #xD | #xA)+
func (p *Parser) parseSpace() error {
	if !isSpace(p.Get()) {
		return fmt.Errorf("expected space")
	}
	p.skipSpace()
	return nil
}

// (#x20 | #x9 | #xD | #xA)*
func (p *Parser) skipSpace() {
	for isSpace(p.Get()) {
		p.Step()
	}
}

// Eq ::= S? '=' S?
func (p *Parser) parseEq() error {
	p.skipSpace()
	if err := p.Must('='); err != nil {
		return err
	}
	p.skipSpace()
	return nil
}

// VersionNum ::= ([a-zA-Z0-9_.:] | '-')+
func (p *Parser) parseVersionNum() (string, error) {
	isVerChar := func() (rune, bool) {
		if isNum(p.Get()) || isAlpha(p.Get()) || p.Test('_') || p.Test('.') || p.Test(':') || p.Test('-') {
			return p.Get(), true
		} else {
			return 0, false
		}
	}

	var str string
	r, ok := isVerChar()
	if !ok {
		return "", fmt.Errorf("error while parsing version number")
	}
	str += string(r)
	p.Step()

	for {
		r, ok := isVerChar()
		if !ok {
			break
		}
		str += string(r)
		p.Step()
	}

	return str, nil
}
