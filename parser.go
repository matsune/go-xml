package xml

import (
	"fmt"
	"strings"
)

type Parser struct {
	*Scanner
}

func NewParser(str string) *Parser {
	return &Parser{
		Scanner: NewScanner(str),
	}
}

/// EBNF for XML 1.0
/// http://www.jelks.nu/XML/xmlebnf.html#NT-VersionInfo

// - Document

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

/// - Prolog

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
		if _, err := p.parseMisc(); err != nil {
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
			if _, err := p.parseMisc(); err != nil {
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
func (p *Parser) parseVersion() (string, error) {
	var err error
	if err = p.parseSpace(); err != nil {
		return "", err
	}

	if err = p.Musts("version"); err != nil {
		return "", err
	}
	if err = p.parseEq(); err != nil {
		return "", err
	}

	var quote rune
	quote, err = p.parseQuote()
	if err != nil {
		return "", err
	}

	var ver string
	ver, err = p.parseVersionNum()
	if err != nil {
		return "", err
	}

	if err = p.Must(quote); err != nil {
		return "", err
	}
	return ver, nil
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

func (p *Parser) isMisc() bool {
	return p.Tests(`<!--`) || p.Tests(`<?`) || isSpace(p.Get())
}

// Misc ::= Comment | PI | S
func (p *Parser) parseMisc() (interface{}, error) {
	if p.Tests(`<!--`) {
		// ignore comment
		c, err := p.parseComment()
		if err != nil {
			return nil, err
		}
		return c, nil
	} else if p.Tests(`<?`) {
		// ignore PI
		pi, err := p.parsePI()
		if err != nil {
			return nil, err
		}
		return pi, err
	} else if isSpace(p.Get()) {
		p.skipSpace()
		return nil, nil
	} else {
		return nil, p.errorf("error while parsing misc")
	}
}

/// - White Space

// S ::= (#x20 | #x9 | #xD | #xA)+
func (p *Parser) parseSpace() error {
	if !isSpace(p.Get()) {
		return p.errorf("expected space")
	}
	p.skipSpace()
	return nil
}

// Skip spaces until reaches not space
func (p *Parser) skipSpace() {
	for isSpace(p.Get()) {
		p.Step()
	}
}

/// - Names and Tokens

// NameChar ::= Letter | Digit | '.' | '-' | '_' | ':' | CombiningChar |  Extender
func (p *Parser) isNameChar() bool {
	return isLetter(p.Get()) || isDigit(p.Get()) || p.Test('.') || p.Test('-') || p.Test('_') || p.Test(':') || isCombining(p.Get()) || isExtender(p.Get())
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

/// - Literals

// EntityValue ::= '"' ([^%&"] | PEReference | Reference)* '"' |  "'" ([^%&'] |  PEReference |  Reference)* "'"
func (p *Parser) parseEntityValue() (EntityValue, error) {
	var quote rune
	var err error
	if quote, err = p.parseQuote(); err != nil {
		return nil, err
	}

	res := EntityValue{}

	var str string
	for {
		if p.Test(quote) || p.isEnd() {
			break
		}

		cur := p.cursor

		if p.Test('&') {
			if len(str) > 0 {
				res = append(res, str)
				str = ""
			}

			// try EntityRef
			var eRef *EntityRef
			eRef, err = p.parseEntityRef()
			if err != nil {
				p.cursor = cur
				// try CharRef
				var cRef *CharRef
				cRef, err = p.parseCharRef()
				if err != nil {
					return nil, p.errorf("error AttValue")
				}
				res = append(res, cRef)
			} else {
				res = append(res, eRef)
			}
		} else if p.Test('%') {
			if len(str) > 0 {
				res = append(res, str)
				str = ""
			}

			var pRef *PERef
			if pRef, err = p.parsePERef(); err != nil {
				return nil, err
			}
			res = append(res, pRef)
		} else {
			str += string(p.Get())
			p.Step()
		}
	}

	if len(str) > 0 {
		res = append(res, str)
		str = ""
	}

	if err = p.Must(quote); err != nil {
		return nil, err
	}

	return res, nil
}

// AttValue ::= '"' ([^<&"] | Reference)* '"' |  "'" ([^<&'] | Reference)* "'"
func (p *Parser) parseAttValue() (AttValue, error) {
	var quote rune
	var err error
	if quote, err = p.parseQuote(); err != nil {
		return nil, err
	}

	res := AttValue{}

	var str string
	for {
		if p.Test('<') {
			return nil, p.errorf("error AttValue")
		}
		if p.Test(quote) || p.isEnd() {
			break
		}

		cur := p.cursor

		if p.Test('&') {
			if len(str) > 0 {
				res = append(res, str)
				str = ""
			}

			// try EntityRef
			var eRef *EntityRef
			eRef, err = p.parseEntityRef()
			if err != nil {
				p.cursor = cur
				// try CharRef
				var cRef *CharRef
				cRef, err = p.parseCharRef()
				if err != nil {
					return nil, p.errorf("error AttValue")
				}
				res = append(res, cRef)
			} else {
				res = append(res, eRef)
			}
		} else {
			str += string(p.Get())
			p.Step()
		}
	}

	if len(str) > 0 {
		res = append(res, str)
		str = ""
	}

	if err = p.Must(quote); err != nil {
		return nil, err
	}

	return res, nil
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

		if p.isEnd() {
			return "", p.errorf("could not find quote %c", quote)
		}
	}
	p.Step()

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

/// - Comments

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

/// - Processing Instructions

// PI ::= ::= '<?' PITarget (S (Char* - (Char* '?>' Char*)))? '?>'
func (p *Parser) parsePI() (*PI, error) {
	var err error
	if err = p.Musts("<?"); err != nil {
		return nil, err
	}
	var pi PI
	if pi.Target, err = p.parsePITarget(); err != nil {
		return nil, err
	}

	if isSpace(p.Get()) {
		p.skipSpace()

		for !p.Tests("?>") && !p.isEnd() && isChar(p.Get()) {
			pi.Instruction += string(p.Get())
			p.Step()
		}
	}

	if err = p.Musts("?>"); err != nil {
		return nil, err
	}
	return &pi, nil
}

// PITarget ::= Name - (('X' | 'x') ('M' | 'm') ('L' | 'l'))
func (p *Parser) parsePITarget() (string, error) {
	var n string
	var err error
	if n, err = p.parseName(); err != nil {
		return "", err
	}
	if strings.ContainsAny(n, "xmlXML") {
		return "", fmt.Errorf("error parsing PITarget")
	}
	return n, nil
}

/// - Document Type Definition

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
				m, err = p.parseMarkup()
				if err != nil {
					return nil, err
				}
				d.Markups = append(d.Markups, m)
			} else if p.Test('%') {
				var ref *PERef
				ref, err = p.parsePERef()
				if err != nil {
					return nil, err
				}
				d.PERef = ref
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

// markupdecl ::= elementdecl |  AttlistDecl |  EntityDecl |  NotationDecl | PI |  Comment
func (p *Parser) parseMarkup() (Markup, error) {
	var err error
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
	default:
		err = p.errorf("error while parsing markup")
	}
	return m, err
}

/// - Standalone Document Declaration

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

/// - Element Type Declaration

// elementdecl ::= '<!ELEMENT' S Name S contentspec S? '>'
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

// contentspec ::= 'EMPTY' | 'ANY' | Mixed | children
func (p *Parser) parseContentSpec() (ContentSpec, error) {
	if p.Tests("EMPTY") {
		p.StepN(len("EMPTY"))
		return &EMPTY{}, nil
	} else if p.Tests("ANY") {
		p.StepN(len("ANY"))
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

/// - Element-content Models

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

// cp ::= (Name | choice | seq) ('?' | '*' | '+')?
func (p *Parser) parseCP() (*CP, error) {
	var cp CP
	var err error
	if p.Test('(') { // choice or seq
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

/// - Mixed-content Declaration

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
			p.Step()
			break
		}
		if p.isEnd() {
			return nil, p.errorf("could not find ')'")
		}
		var err error
		if err = p.Must('|'); err != nil {
			return nil, err
		}
		p.skipSpace()
		var n string
		n, err = p.parseName()
		if err != nil {
			return nil, err
		}
		m.Names = append(m.Names, n)
	}

	if len(m.Names) > 0 {
		if err := p.Must('*'); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

/// - Attribute-list Declaration

// AttlistDecl ::= '<!ATTLIST' S Name AttDef* S? '>'
func (p *Parser) parseAttlist() (*Attlist, error) {
	var att Attlist
	var err error
	if err = p.Musts("<!ATTLIST"); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}
	if att.Name, err = p.parseName(); err != nil {
		return nil, err
	}
	// S Name  or  S? '>'
	for {
		cur := p.cursor

		p.skipSpace()
		if p.Test('>') {
			p.cursor = cur
			break
		}
		if p.isEnd() {
			return nil, p.errorf("error while parsing Attlist")
		}
		p.cursor = cur

		var def *AttDef
		if def, err = p.parseAttDef(); err != nil {
			return nil, err
		}
		att.Defs = append(att.Defs, def)
	}

	p.skipSpace()
	if err = p.Must('>'); err != nil {
		return nil, err
	}

	return &att, nil
}

// AttDef ::= S Name S AttType S DefaultDecl
func (p *Parser) parseAttDef() (*AttDef, error) {
	var err error
	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	var def AttDef
	if def.Name, err = p.parseName(); err != nil {
		return nil, err
	}

	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	if def.Type, err = p.parseAttType(); err != nil {
		return nil, err
	}

	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	if def.Decl, err = p.parseDefaultDecl(); err != nil {
		return nil, err
	}

	return &def, nil
}

/// - Attribute Types

// AttType ::= StringType | TokenizedType | EnumeratedType
func (p *Parser) parseAttType() (AttType, error) {
	if p.Tests(string(Att_CDATA)) {
		p.StepN(len(Att_CDATA))
		return Att_CDATA, nil
	} else if p.Tests(string(Att_ID)) || p.Tests(string(Att_IDREF)) || p.Tests(string(Att_IDREFS)) || p.Tests(string(Att_ENTITY)) || p.Tests(string(Att_ENTITIES)) || p.Tests(string(Att_NMTOKEN)) || p.Tests(string(Att_NMTOKENS)) {
		var tok TokenizedType
		switch {
		case p.Tests(string(Att_ID)):
			tok = Att_ID
			if p.Tests(string(Att_IDREF)) {
				tok = Att_IDREF
				if p.Tests(string(Att_IDREFS)) {
					tok = Att_IDREFS
				}
			}
		case p.Tests(string(Att_ENTITY)):
			tok = Att_ENTITY
		case p.Tests(string(Att_ENTITIES)):
			tok = Att_ENTITIES
		case p.Tests(string(Att_NMTOKEN)):
			tok = Att_NMTOKEN
			if p.Tests(string(Att_NMTOKENS)) {
				tok = Att_NMTOKENS
			}
		}
		p.StepN(len(tok))
		return tok, nil
	} else if p.Tests("NOTATION") {
		return p.parseNotationType()
	} else if p.Test('(') {
		return p.parseEnum()
	}
	return nil, p.errorf("error while parsing AttType")
}

/// - Enumerated Attribute Types

// NotationType ::= 'NOTATION' S '(' S? Name (S? '|' S? Name)* S? ')'
func (p *Parser) parseNotationType() (*NotationType, error) {
	var err error
	if err = p.Musts("NOTATION"); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}
	if err = p.Must('('); err != nil {
		return nil, err
	}
	p.skipSpace()

	var n NotationType
	var t string
	if t, err = p.parseName(); err != nil {
		return nil, err
	}
	n.Names = append(n.Names, t)

	for {
		cur := p.cursor

		p.skipSpace()
		if p.Test(')') {
			p.Step()
			break
		}
		if p.isEnd() {
			return nil, p.errorf("could not found ')'")
		}
		p.cursor = cur

		p.skipSpace()
		if err = p.Must('|'); err != nil {
			return nil, err
		}
		p.skipSpace()

		if t, err = p.parseName(); err != nil {
			return nil, err
		}
		n.Names = append(n.Names, t)
	}

	return &n, nil
}

// Enumeration ::= '(' S? Nmtoken (S? '|' S? Nmtoken)* S? ')'
func (p *Parser) parseEnum() (*Enum, error) {
	var err error
	var e Enum

	if err = p.Must('('); err != nil {
		return nil, err
	}
	p.skipSpace()

	var nm string
	if nm, err = p.parseNmtoken(); err != nil {
		return nil, err
	}
	e.Cases = append(e.Cases, nm)

	for {
		cur := p.cursor

		p.skipSpace()
		if p.Test(')') {
			p.Step()
			break
		}
		if p.isEnd() {
			return nil, p.errorf("could not found ')'")
		}
		p.cursor = cur

		p.skipSpace()
		if err = p.Must('|'); err != nil {
			return nil, err
		}
		p.skipSpace()

		if nm, err = p.parseNmtoken(); err != nil {
			return nil, err
		}
		e.Cases = append(e.Cases, nm)
	}

	return &e, nil
}

// Nmtoken ::= (NameChar)+
func (p *Parser) parseNmtoken() (string, error) {
	var str string
	r := p.Get()
	for isNameChar(r) {
		str += string(r)
		p.Step()
		r = p.Get()
	}
	if len(str) == 0 {
		return "", p.errorf("error Nmtoken")
	}
	return str, nil
}

/// - Attribute Defaults

// DefaultDecl ::= '#REQUIRED' | '#IMPLIED' | (('#FIXED' S)? AttValue)
func (p *Parser) parseDefaultDecl() (*DefaultDecl, error) {
	var d DefaultDecl
	var err error
	if p.Tests(string(REQUIRED)) {
		p.StepN(len(REQUIRED))
		d.Type = REQUIRED
		return &d, nil
	} else if p.Tests(string(IMPLIED)) {
		p.StepN(len(IMPLIED))
		d.Type = IMPLIED
		return &d, nil
	} else {
		if p.Tests(string(FIXED)) {
			p.StepN(len(FIXED))
			if err = p.parseSpace(); err != nil {
				return nil, err
			}
		}
		d.Type = FIXED
		if d.AttValue, err = p.parseAttValue(); err != nil {
			return nil, err
		}
		return &d, nil
	}
}

/// - Character Reference

// CharRef ::= '&#' [0-9]+ ';' | '&#x' [0-9a-fA-F]+ ';'
func (p *Parser) parseCharRef() (*CharRef, error) {
	var ref CharRef
	var err error

	if p.Tests("&#x") {
		ref.Prefix = "&#x"
		p.StepN(len("&#x"))

		r := p.Get()
		if !isNum(r) && !isAlpha(r) {
			return nil, p.errorf("error CharRef")
		}

		for isNum(r) || isAlpha(r) {
			ref.Value += string(r)
			p.Step()
			r = p.Get()
		}
	} else if p.Tests("&#") {
		ref.Prefix = "&#"
		p.StepN(len("&#"))

		r := p.Get()
		if !isNum(r) {
			return nil, p.errorf("error CharRef")
		}

		for isNum(r) {
			ref.Value += string(r)
			p.Step()
			r = p.Get()
		}
	} else {
		return nil, p.errorf("error CharRef")
	}

	if err = p.Must(';'); err != nil {
		return nil, err
	}

	return &ref, nil
}

/// - Entity Reference

// EntityRef ::= '&' Name ';'
func (p *Parser) parseEntityRef() (*EntityRef, error) {
	var err error
	if err = p.Must('&'); err != nil {
		return nil, err
	}
	var e EntityRef
	e.Name, err = p.parseName()
	if err != nil {
		return nil, err
	}
	if err = p.Must(';'); err != nil {
		return nil, err
	}
	return &e, nil
}

// PEReference ::= '%' Name ';'
func (p *Parser) parsePERef() (*PERef, error) {
	var err error
	if err = p.Must('%'); err != nil {
		return nil, err
	}
	var e PERef
	e.Name, err = p.parseName()
	if err != nil {
		return nil, err
	}
	if err = p.Must(';'); err != nil {
		return nil, err
	}
	return &e, nil
}

/// - Entity Declaration

// EntityDecl ::= GEDecl |  PEDecl
func (p *Parser) parseEntity() (*Entity, error) {
	var err error
	if err = p.Musts("<!ENTITY"); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	var e Entity

	// PEDecl ::= '<!ENTITY' S '%' S Name S PEDef S? '>'
	// GEDecl ::= '<!ENTITY' S Name S EntityDef S? '>'

	if p.Test('%') {
		e.Type = EntityType_PE

		p.Step()

		if err = p.parseSpace(); err != nil {
			return nil, err
		}
	} else {
		e.Type = EntityType_GE
	}

	if e.Name, err = p.parseName(); err != nil {
		return nil, err
	}

	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	if e.Type == EntityType_PE {
		// PEDef
		e.Value, e.ExtID, err = p.parsePEDef()
	} else {
		// EntityDef
		e.Value, e.ExtID, e.NData, err = p.parseEntityDef()
	}
	if err != nil {
		return nil, err
	}

	p.skipSpace()

	if err = p.Must('>'); err != nil {
		return nil, err
	}

	return &e, nil
}

// EntityDef ::= EntityValue | (ExternalID NDataDecl?)
func (p *Parser) parseEntityDef() (EntityValue, *ExternalID, string, error) {
	var value EntityValue
	var ndata string
	var ext *ExternalID
	var err error

	if p.Test('\'') || p.Test('"') {
		value, err = p.parseEntityValue()
		if err != nil {
			return nil, nil, "", err
		}
		return value, nil, "", nil
	} else if p.Tests("SYSTEM") || p.Tests("PUBLIC") {
		ext, err = p.parseExternalID()
		if err != nil {
			return nil, nil, "", err
		}

		cur := p.cursor

		ndata, err = p.parseNData()
		if err != nil {
			p.cursor = cur
		}

		return nil, ext, ndata, nil
	} else {
		return nil, nil, "", p.errorf("error EntityDef")
	}
}

// PEDef ::= EntityValue | ExternalID
func (p *Parser) parsePEDef() (EntityValue, *ExternalID, error) {
	var value EntityValue
	var ext *ExternalID
	var err error

	if p.Test('\'') || p.Test('"') {
		value, err = p.parseEntityValue()
		if err != nil {
			return nil, nil, err
		}
		return value, nil, nil
	} else if p.Tests("SYSTEM") || p.Tests("PUBLIC") {
		ext, err = p.parseExternalID()
		if err != nil {
			return nil, nil, err
		}

		return nil, ext, nil
	} else {
		return nil, nil, p.errorf("error EntityDef")
	}
}

/// - External Entity Declaration

// ExternalID ::= 'SYSTEM' S SystemLiteral | 'PUBLIC' S PubidLiteral S SystemLiteral
func (p *Parser) parseExternalID() (*ExternalID, error) {
	var ext ExternalID
	if p.Tests(string(ExtSystem)) {
		p.StepN(len(ExtSystem))
		ext.Identifier = ExtSystem
	} else if p.Tests(string(ExtPublic)) {
		p.StepN(len(ExtPublic))
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

// NDataDecl ::= S 'NDATA' S Name
func (p *Parser) parseNData() (string, error) {
	var err error
	if err = p.parseSpace(); err != nil {
		return "", err
	}
	if err = p.Musts("NDATA"); err != nil {
		return "", err
	}
	if err = p.parseSpace(); err != nil {
		return "", err
	}
	return p.parseName()
}

/// - Encoding Declaration

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

/// - Notation Declarations

// NotationDecl ::= '<!NOTATION' S Name S (ExternalID | PublicID) S? '>'
func (p *Parser) parseNotation() (*Notation, error) {
	var n Notation
	var err error
	if err = p.Musts("<!NOTATION"); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}
	if n.Name, err = p.parseName(); err != nil {
		return nil, err
	}
	if err = p.parseSpace(); err != nil {
		return nil, err
	}

	// ExternalID ::= 'SYSTEM' S SystemLiteral | 'PUBLIC' S PubidLiteral
	var ext ExternalID
	if p.Tests(string(ExtSystem)) {
		p.StepN(len(ExtSystem))
		ext.Identifier = ExtSystem
	} else if p.Tests(string(ExtPublic)) {
		p.StepN(len(ExtPublic))
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

		// ( S SystemLiteral )?
		cur := p.cursor

		err = p.parseSpace()
		if err == nil {
			var sys string
			sys, err = p.parseSystemLiteral()
			if err == nil {
				ext.System = sys
			} else {
				err = nil
				p.cursor = cur
			}
		} else {
			p.cursor = cur
			err = nil
		}
	} else {
		var sys string
		sys, err = p.parseSystemLiteral()
		if err != nil {
			return nil, err
		}
		ext.System = sys
	}

	n.ExtID = ext

	p.skipSpace()

	if err = p.Must('>'); err != nil {
		return nil, err
	}

	return &n, nil
}

/// - Others

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
