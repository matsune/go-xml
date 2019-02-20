package xml

import "fmt"

type Parser struct {
	*Scanner
}

func NewParser(str string) *Parser {
	return &Parser{
		Scanner: NewScanner(str),
	}
}

func unimplemented(method string) {
	panic("Unimplemented " + method)
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

	p.skipSpace()

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

	if p.Tests(`<!DOCTYPE`) {
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

	// keep cursor at this time
	cur := p.cursor

	p.skipSpace()

	if p.Tests("encoding") {
		// reset cursor before skipping spaces
		p.cursor = cur

		enc, err := p.parseEncoding()
		if err != nil {
			return nil, err
		}
		x.Encoding = enc
	} else {
		p.cursor = cur
	}

	cur = p.cursor

	p.skipSpace()

	if p.Tests("standalone") {
		p.cursor = cur

		std, err := p.parseStandalone()
		if err != nil {
			return nil, err
		}
		x.Standalone = std
	} else {
		p.cursor = cur
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
		err = p.errorf("expected ' or \"")
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
		return "", p.errorf("error while parsing encoding name")
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
		return false, p.errorf("error while parsing standalone")
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
		// ignore PI
		_, err := p.parsePI()
		if err != nil {
			return err
		}
	} else if isSpace(p.Get()) {
		p.skipSpace()
	} else {
		return p.errorf("error while parsing misc")
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
			return "", p.errorf("error while parsing comment")
		}
	}

	if err := p.Musts(`-->`); err != nil {
		return "", err
	}

	return str, nil
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
		p.Step()

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
	}

	p.skipSpace()

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
		return nil, p.errorf("error while parsing ExternalID")
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

	return &ext, nil
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

	if err = p.Must(quote); err != nil {
		return "", err
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
		return "", p.errorf("error while parsing name")
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
		return p.errorf("expected space")
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
		return "", p.errorf("error while parsing version number")
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
	unimplemented("parseAttlist")
	return nil, nil
}

func (p *Parser) parseEntity() (*Entity, error) {
	unimplemented("parseEntity")
	return nil, nil
}

func (p *Parser) parseNotation() (*Notation, error) {
	unimplemented("parseNotation")
	return nil, nil
}

func (p *Parser) parsePI() (*PI, error) {
	unimplemented("parsePI")
	return nil, nil
}

func (p *Parser) parsePEReference() (string, error) {
	unimplemented("parsePEReference")
	return "", nil
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

		cur := p.cursor
		{ // try parsing mixed
			var m *Mixed
			m, err = p.parseMixed()
			if err == nil {
				return m, nil
			}
		}
		// reset cursor if it wasn't mixed
		p.cursor = cur

		var ch *Children
		ch, err = p.parseChildren()
		if err != nil {
			return nil, err
		}
		return ch, nil
	}
}

// children ::= (choice | seq) ('?' | '*' | '+')?
func (p *Parser) parseChildren() (*Children, error) {
	var c Children
	var err error

	cur := p.cursor
	{
		var choice *Choice
		choice, err = p.parseChoice()
		if err == nil {
			c.ChoiceSeq = choice
			if p.Test('?') || p.Test('*') || p.Test('+') {
				r := p.Get()
				c.Suffix = &r
				p.Step()
			}
			return &c, nil
		}
	}
	p.cursor = cur

	var s *Seq
	s, err = p.parseSeq()
	if err == nil {
		c.ChoiceSeq = s
		if p.Test('?') || p.Test('*') || p.Test('+') {
			r := p.Get()
			c.Suffix = &r
			p.Step()
		}
		return &c, nil
	}

	return nil, p.errorf("error while parsing children")
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

	if len(m.Names) > 0 {
		if err := p.Must('*'); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

// cp ::= (Name | choice | seq) ('?' | '*' | '+')?
func (p *Parser) parseCP() (*CP, error) {
	var cp CP
	var err error
	if p.Test('(') {
		// choice or seq
		cur := p.cursor

		var choice *Choice
		choice, err = p.parseChoice()
		if err != nil {
			p.cursor = cur

			var seq *Seq
			seq, err = p.parseSeq()
			if err != nil {
				return nil, err
			}
			cp.ChoiceSeq = seq
		} else {
			cp.ChoiceSeq = choice
		}
	} else {
		var n string
		n, err = p.parseName()
		if err != nil {
			return nil, err
		}
		cp.Name = n
	}

	if p.Test('?') || p.Test('*') || p.Test('+') {
		r := p.Get()
		cp.Suffix = &r
		p.Step()
	}

	return &cp, nil
}

// choice ::= '(' S? cp ( S? '|' S? cp )* S? ')'
func (p *Parser) parseChoice() (*Choice, error) {
	if err := p.Must('('); err != nil {
		return nil, err
	}
	p.skipSpace()
	var cps []CP
	cp, err := p.parseCP()
	if err != nil {
		return nil, err
	}
	cps = append(cps, *cp)
	for {
		p.skipSpace()
		if !p.Test('|') {
			break
		}
		p.Step()

		cp, err = p.parseCP()
		if err != nil {
			return nil, err
		}
		cps = append(cps, *cp)
	}
	if err := p.Must(')'); err != nil {
		return nil, err
	}
	return &Choice{
		CPs: cps,
	}, nil
}

// seq ::= '(' S? cp ( S? ',' S? cp )* S? ')'
func (p *Parser) parseSeq() (*Seq, error) {
	if err := p.Must('('); err != nil {
		return nil, err
	}
	p.skipSpace()
	var cps []CP
	cp, err := p.parseCP()
	if err != nil {
		return nil, err
	}

	cps = append(cps, *cp)
	for {
		p.skipSpace()
		if !p.Test(',') {
			break
		}
		p.Step()

		cp, err = p.parseCP()
		if err != nil {
			return nil, err
		}
		cps = append(cps, *cp)
	}

	if err := p.Must(')'); err != nil {
		return nil, err
	}
	return &Seq{
		CPs: cps,
	}, nil
}
