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

func (p *Parser) removeLast() {
	if len(p.stack) > 0 {
		p.stack = p.stack[:len(p.stack)-1]
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

// document ::= prolog element Misc*
func (p *Parser) Parse() error {
	pro, err := p.parseProlog()
	if err != nil {
		return err
	}
	fmt.Println(">>>", pro)
	// TODO: parse element
	return nil
}

// prolog ::= XMLDecl? Misc* (doctypedecl Misc*)?
func (p *Parser) parseProlog() (*Prolog, error) {
	pro := Prolog{}
	if p.Tests("<?xml") {
		xmlDecl, err := p.parseXmlDecl()
		if err != nil {
			return nil, err
		}
		pro.XMLDecl = xmlDecl
	}
	for p.isMisc() {
		if err := p.parseMisc(); err != nil {
			return nil, err
		}
	}

	if p.isDoctype() {
		doc, err := p.parseDoctype()
		if err != nil {
			return nil, err
		}
		pro.DOCType = doc

		for p.isMisc() {
			if err := p.parseMisc(); err != nil {
				return nil, err
			}
		}
	}
	return &pro, nil
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
	x.Version = ver

	p.push()
	p.skipSpace()
	if p.Tests("encoding") {
		p.pop()

		enc, err := p.parseEncoding()
		if err != nil {
			return nil, err
		}
		x.Encoding = enc
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

	var quote rune
	quote, err = p.parseQuote()
	if err != nil {
		return
	}

	ver, err = p.parseVersionNum()
	if err != nil {
		return
	}

	if err = p.Must(quote); err != nil {
		return
	}

	return
}

func (p *Parser) parseQuote() (rune, error) {
	var err error
	r := p.Get()
	if isQuote(r) {
		p.Step()
	} else {
		err = fmt.Errorf("expected ' or \"")
	}
	return r, err
}

// EncodingDecl ::= S 'encoding' Eq ('"' EncName  '"' |  "'" EncName "'" )
func (p *Parser) parseEncoding() (string, error) {
	var err error
	if err = p.parseSpace(); err != nil {
		return "", err
	}
	if err = p.Musts("encoding"); err != nil {
		return "", err
	}
	if err = p.parseEq(); err != nil {
		return "", err
	}

	var quote rune
	if quote, err = p.parseQuote(); err != nil {
		return "", err
	}

	var enc string
	enc, err = p.parseEncName()
	if err != nil {
		return "", err
	}

	if err = p.Must(quote); err != nil {
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
	var err error
	if err = p.parseSpace(); err != nil {
		return false, err
	}
	if err = p.Musts("standalone"); err != nil {
		return false, err
	}
	if err = p.parseEq(); err != nil {
		return false, err
	}
	var quote rune
	if quote, err = p.parseQuote(); err != nil {
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
	if err = p.Must(quote); err != nil {
		return false, err
	}
	return std, nil
}

func (p *Parser) isMisc() bool {
	return p.Tests(`<!--`) || p.Tests(`<?`) || isSpace(p.Get())
}

// Misc ::= Comment | PI | S
func (p *Parser) parseMisc() error {
	if p.Tests(`<!--`) {
		// ignore comment
		_, err := p.parseComment()
		if err != nil {
			return err
		}
	} else if p.Tests(`<?`) {
		panic("unimplemented PI")
	} else if isSpace(p.Get()) {
		p.skipSpace()
	} else {
		return fmt.Errorf("error while parsing misc")
	}
	return nil
}

// Comment ::= '<!--' ((Char - '-') | ('-' (Char - '-')))* '-->'
func (p *Parser) parseComment() (Comment, error) {
	if err := p.Musts(`<!--`); err != nil {
		return "", err
	}

	var str Comment
	for !p.Tests("--") {
		r := p.Get()
		if isChar(r) {
			str += Comment(r)
			p.Step()
		} else {
			return "", fmt.Errorf("error while parsing comment")
		}
	}

	if err := p.Musts(`-->`); err != nil {
		return "", err
	}

	return str, nil
}

func (p *Parser) isDoctype() bool {
	return p.Tests(`<!DOCTYPE`)
}

// doctypedecl ::= '<!DOCTYPE' S Name (S ExternalID)? S? ('[' (markupdecl | PEReference | S)* ']' S?)? '>'
func (p *Parser) parseDoctype() (*DOCType, error) {
	if err := p.Musts(`<!DOCTYPE`); err != nil {
		return nil, err
	}
	if err := p.parseSpace(); err != nil {
		return nil, err
	}
	var d DOCType
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	d.Name = name

	p.skipSpace()

	if p.Tests(string(ExtSystem)) || p.Tests(string(ExtPublic)) {
		var ext *ExternalID
		ext, err = p.parseExternalID()
		if err != nil {
			return nil, err
		}
		d.ExternalID = ext
	} else if p.Test('[') {
		err = p.Must('[')
		if err != nil {
			return nil, err
		}
		for {
			if p.Tests("<!ELEMENT") || p.Tests("<!ATTLIST") || p.Tests("<!ENTITY") || p.Tests("<!NOTATION") || p.Tests("<?") || p.Tests("<!--") {
				var m Markup
				switch {
				case p.Tests("<!ELEMENT"):
					m, err = p.parseElement()
				case p.Tests("<!ATTLIST"):
					m, err = p.parseAttlist()
				case p.Tests("<!ENTITY"):
					m, err = p.parseEntity()
				case p.Tests("<!NOTATION"):
					m, err = p.parseNotation()
				case p.Tests("<?"):
					m, err = p.parsePI()
				case p.Tests("<!--"):
					m, err = p.parseComment()
				}
				if err != nil {
					return nil, err
				}
				d.Markups = append(d.Markups, m)
			} else if p.Test('%') {
				var ref string
				ref, err = p.parsePEReference()
				if err != nil {
					return nil, err
				}
				d.PEReference = ref
			} else if isSpace(p.Get()) {
				err = p.parseSpace()
				if err != nil {
					return nil, err
				}
			} else {
				break
			}
		}
		err = p.Must(']')
		if err != nil {
			return nil, err
		}
		p.skipSpace()
	}

	err = p.Must('>')
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ExternalID ::= 'SYSTEM' S SystemLiteral | 'PUBLIC' S PubidLiteral S SystemLiteral
func (p *Parser) parseExternalID() (*ExternalID, error) {
	var ext ExternalID
	if p.Tests(string(ExtSystem)) {
		if err := p.Musts(string(ExtSystem)); err != nil {
			return nil, err
		}
		ext.Identifier = ExtSystem
	} else if p.Tests(string(ExtPublic)) {
		if err := p.Musts(string(ExtPublic)); err != nil {
			return nil, err
		}
		ext.Identifier = ExtPublic
	} else {
		return nil, fmt.Errorf("error while parsing ExternalID")
	}

	if err := p.parseSpace(); err != nil {
		return nil, err
	}

	if ext.Identifier == ExtPublic {
		pubid, err := p.parsePubidLiteral()
		if err != nil {
			return nil, err
		}
		ext.Pubid = pubid

		if err := p.parseSpace(); err != nil {
			return nil, err
		}
	}

	sys, err := p.parseSystemLiteral()
	if err != nil {
		return nil, err
	}
	ext.System = sys

	return nil, nil
}

// SystemLiteral ::= ('"' [^"]* '"') | ("'" [^']* "'")
func (p *Parser) parseSystemLiteral() (string, error) {
	var quote rune
	var err error
	if quote, err = p.parseQuote(); err != nil {
		return "", err
	}

	var lit string
	for !p.Test(quote) {
		lit += string(p.Get())
		p.Step()
	}

	return lit, nil
}

// PubidLiteral ::= '"' PubidChar* '"' | "'" (PubidChar - "'")* "'"
func (p *Parser) parsePubidLiteral() (string, error) {
	var quote rune
	var err error
	if quote, err = p.parseQuote(); err != nil {
		return "", err
	}

	var lit string
	r := p.Get()
	for isPubidChar(r) {
		if r == '\'' && quote == '\'' {
			break
		}
		lit += string(r)
		p.Step()
		r = p.Get()
	}

	if err = p.Must(quote); err != nil {
		return "", err
	}

	return lit, nil
}

// Name ::= (Letter | '_' | ':') (NameChar)*
func (p *Parser) parseName() (string, error) {
	var n string
	if isLetter(p.Get()) || p.Test('_') || p.Test(':') {
		n += string(p.Get())
		p.Step()
	} else {
		return "", fmt.Errorf("error while parsing name")
	}
	for p.isNameChar() {
		n += string(p.Get())
		p.Step()
	}
	return n, nil
}

func (p *Parser) isNameChar() bool {
	return isLetter(p.Get()) || isDigit(p.Get()) || p.Test('.') || p.Test('-') || p.Test('_') || p.Test(':') || isCombining(p.Get()) || isExtender(p.Get())
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

func (p *Parser) parseElement() (*Element, error) {
	var err error
	if err = p.Musts("<!ELEMENT"); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}
	var n string
	if n, err = p.parseName(); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}
	var c ContentSpec
	if c, err = p.parseContentSpec(); err != nil {
		return nil, err
	}
	p.skipSpace()
	if err = p.Must('>'); err != nil {
		return nil, err
	}
	return &Element{
		Name:        n,
		ContentSpec: c,
	}, nil
}

func (p *Parser) parseAttlist() (*Attlist, error) {
	panic("unimplemented parseAttlist")
}

func (p *Parser) parseEntity() (*Entity, error) {
	panic("unimplemented parseEntity")
}

func (p *Parser) parseNotation() (*Notation, error) {
	panic("unimplemented parseNotation")
}

func (p *Parser) parsePI() (*PI, error) {
	panic("unimplemented parsePI")
}

func (p *Parser) parsePEReference() (string, error) {
	panic("unimplemented parsePEReference")
}

// contentspec ::= 'EMPTY' | 'ANY' | Mixed | children
func (p *Parser) parseContentSpec() (ContentSpec, error) {
	if p.Tests("EMPTY") {
		if err := p.Musts("EMPTY"); err != nil {
			return nil, err
		}
		return &EMPTY{}, nil
	} else if p.Tests("ANY") {
		if err := p.Musts("ANY"); err != nil {
			return nil, err
		}
		return &ANY{}, nil
	} else {
		var err error

		p.push()
		{ // try parsing mixed
			var m *Mixed
			m, err = p.parseMixed()
			if err == nil {
				p.removeLast()
				return m, nil
			}
		}
		p.pop()

		p.push()
		{ // try parsing choice
			var c *Choice
			c, err = p.parseChoice()
			if err == nil {
				p.removeLast()
				return c, nil
			}
		}
		p.pop()

		p.push()
		{
			var s *Seq
			s, err = p.parseSeq()
			if err == nil {
				p.removeLast()
				return s, nil
			}
		}
		p.pop()

		return nil, err
	}
}

// Mixed ::= '(' S? '#PCDATA' (S? '|' S? Name)* S? ')*' | '(' S? '#PCDATA' S? ')'
func (p *Parser) parseMixed() (*Mixed, error) {
	if err := p.Must('('); err != nil {
		return nil, err
	}
	p.skipSpace()
	if err := p.Musts("#PCDATA"); err != nil {
		return nil, err
	}

	var m Mixed
	for {
		p.skipSpace()
		if p.Test(')') {
			break
		}
		var err error
		if err = p.Must('|'); err != nil {
			return nil, err
		}
		p.skipSpace()
		var n string
		n, err = p.parseName()
		if err != nil {
			return nil, nil
		}
		m.Names = append(m.Names, n)
	}

	if err := p.Must(')'); err != nil {
		return nil, err
	}

	return &m, nil
}

func (p *Parser) parseChoice() (*Choice, error) {
	panic("unimplemented parseMixed")
}

func (p *Parser) parseSeq() (*Seq, error) {
	panic("unimplemented parseSeq")
}
