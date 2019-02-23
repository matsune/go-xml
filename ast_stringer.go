package xml

import (
	"fmt"
	"strings"
)

func (x XMLDecl) String() string {
	str := "<?xml"

	str += fmt.Sprintf(` version="%s"`, x.Version)

	if len(x.Encoding) > 0 {
		str += fmt.Sprintf(` encoding="%s"`, x.Encoding)
	}

	stdStr := "no"
	if x.Standalone {
		stdStr = "yes"
	}
	str += fmt.Sprintf(` standalone="%s" ?>`, stdStr)

	return str
}

func (d DOCType) String() string {
	str := fmt.Sprintf(`<!DOCTYPE %s`, d.Name)

	if d.ExtID != nil {
		str += fmt.Sprintf(" %s", d.ExtID)
	}

	if len(d.Markups) > 0 || d.PERef != nil {
		str += " ["
	}
	for i, m := range d.Markups {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf(`%s`, m)
	}
	if d.PERef != nil {
		str += fmt.Sprintf(` %s`, d.PERef)
	}
	if len(d.Markups) > 0 || d.PERef != nil {
		str += "]"
	}

	str += ">"
	return str
}

func (e ExternalType) String() string {
	if e == EXT_PUBLIC {
		return "PUBLIC"
	} else {
		return "SYSTEM"
	}
}

func (e ExternalID) String() string {
	str := e.Type.String()
	if len(e.Pubid) > 0 {
		str += fmt.Sprintf(" %q", e.Pubid)
	}
	if len(e.System) > 0 {
		str += fmt.Sprintf(" %q", e.System)
	}
	return str
}

func (e ElementDecl) String() string {
	return fmt.Sprintf(`<!ELEMENT %s %s>`, e.Name, e.ContentSpec)
}

func (a Attlist) String() string {
	str := fmt.Sprintf(`<!ATTLIST %s`, a.Name)
	for _, v := range a.Defs {
		str += fmt.Sprintf("%s", v)
	}
	str += ">"
	return str
}

func (e Entity) String() string {
	str := `<!ENTITY`
	if e.Type == ENTITY_PE {
		str += " %"
	}
	str += fmt.Sprintf(" %s", e.Name)
	if len(e.Value) > 0 {
		str += fmt.Sprintf(" %s", e.Value)
	} else {
		str += fmt.Sprintf(" %s", e.ExtID)

		if len(e.NData) > 0 {
			str += fmt.Sprintf(" NDATA %s", e.NData)
		}
	}
	str += ">"
	return str
}

func (n Notation) String() string {
	return fmt.Sprintf(`<!NOTATION %s %s>`, n.Name, n.ExtID)
}

func (p PI) String() string {
	str := fmt.Sprintf(`<?%s`, p.Target)
	if len(p.Instruction) > 0 {
		str += fmt.Sprintf(` %s`, p.Instruction)
	}
	str += "?>"
	return str
}

func (c Comment) String() string {
	return fmt.Sprintf("<!--%s-->", string(c))
}

func (a AttDef) String() string {
	return fmt.Sprintf(" %s %s %s", a.Name, a.Type, a.Decl)
}

func (d DefaultDecl) String() string {
	str := d.Type.String()
	if len(d.AttValue) > 0 {
		str += fmt.Sprintf(" %s", d.AttValue)
	}
	return str
}

func (n NotationType) String() string {
	str := `NOTATION (`
	for i, n := range n.Names {
		if i > 0 {
			str += "|"
		}
		str += n
	}
	str += ")"
	return str
}

func (e Enum) String() string {
	return "(" + strings.Join(e.Cases, "|") + ")"
}

func (d DefaultDeclType) String() string {
	switch d {
	case DECL_REQUIRED:
		return "#REQUIRED"
	case DECL_IMPLIED:
		return "#IMPLIED"
	case DECL_FIXED:
		return "#FIXED"
	default:
		return ""
	}
}

func (a AttToken) String() string {
	switch a {
	case ATT_CDATA:
		return "CDATA"
	case ATT_ID:
		return "ID"
	case ATT_IDREF:
		return "IDREF"
	case ATT_IDREFS:
		return "IDREFS"
	case ATT_ENTITY:
		return "ENTITY"
	case ATT_ENTITIES:
		return "ENTITIES"
	case ATT_NMTOKEN:
		return "NMTOKEN"
	case ATT_NMTOKENS:
		return "NMTOKENS"
	default:
		return ""
	}
}

func (EMPTY) String() string {
	return "EMPTY"
}

func (ANY) String() string {
	return "ANY"
}

func (m Mixed) String() string {
	str := `(#PCDATA`
	for _, n := range m.Names {
		str += "|" + n
	}
	str += ")"
	return str
}

func (c Children) String() string {
	str := fmt.Sprint(c.ChoiceSeq)
	if c.Suffix != nil {
		str += string(*c.Suffix)
	}
	return str
}

func (c CP) String() string {
	var str string
	if c.ChoiceSeq != nil {
		str = fmt.Sprint(c.ChoiceSeq)
	} else {
		str = c.Name
	}
	if c.Suffix != nil {
		str += string(*c.Suffix)
	}
	return str
}

func (c Choice) String() string {
	str := "("
	for i, cp := range c.CPs {
		if i > 0 {
			str += "|"
		}
		str += fmt.Sprint(cp)
	}
	str += ")"
	return str
}

func (s Seq) String() string {
	str := "("
	for i, cp := range s.CPs {
		if i > 0 {
			str += ","
		}
		str += fmt.Sprint(cp)
	}
	str += ")"
	return str
}

func (e EntityValue) String() string {
	str := `"`
	for _, v := range e {
		str += fmt.Sprint(v)
	}
	str += `"`
	return str
}

func (a AttValue) String() string {
	str := `"`
	for _, v := range a {
		str += fmt.Sprint(v)
	}
	str += `"`
	return str
}

func (e CharRef) String() string {
	return fmt.Sprintf("%s%s;", e.Prefix, e.Value)
}
func (e EntityRef) String() string {
	return fmt.Sprintf("&%s;", e.Name)
}
func (e PERef) String() string {
	return "%" + fmt.Sprintf(`%s;`, e.Name)
}

func (e Element) String() string {
	str := fmt.Sprintf(`<%s`, e.Name)
	for _, attr := range e.Attrs {
		str += fmt.Sprintf(` %s`, attr)
	}

	if e.IsEmptyTag {
		str += "/>"
		return str
	}

	str += ">"
	for _, v := range e.Contents {
		str += fmt.Sprint(v)
	}

	str += fmt.Sprintf(`</%s>`, e.Name)
	return str
}

func (a Attribute) String() string {
	return fmt.Sprintf("%s=%s", a.Name, a.AttValue)
}

func (a Attributes) String() string {
	attrs := make([]string, len(a))
	for i, v := range a {
		attrs[i] = v.String()
	}
	return strings.Join(attrs, " ")
}

func (c CData) String() string {
	return fmt.Sprintf("<![CDATA[%s]]>", string(c))
}
