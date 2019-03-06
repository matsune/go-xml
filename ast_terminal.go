package xml

import (
	"fmt"
	"strings"
)

func (e ExternalType) ToString() string {
	if e == ExternalTypePublic {
		return "PUBLIC"
	} else {
		return "SYSTEM"
	}
}

func (e ExternalID) ToString() string {
	str := e.Type.ToString()
	if len(e.Pubid) > 0 {
		str += fmt.Sprintf(" %q", e.Pubid)
	}
	if len(e.System) > 0 {
		str += fmt.Sprintf(" %q", e.System)
	}
	return str
}

func (e ElementDecl) ToString() string {
	return fmt.Sprintf(`<!ELEMENT %s %s>`, e.Name, e.ContentSpec.ToString())
}

func (a Attlist) ToString() string {
	str := fmt.Sprintf(`<!ATTLIST %s`, a.Name)
	for _, v := range a.Defs {
		str += fmt.Sprintf("%s", v.ToString())
	}
	str += ">"
	return str
}

func (e Entity) ToString() string {
	str := `<!ENTITY`
	if e.Type == EntityTypePE {
		str += " %"
	}
	str += fmt.Sprintf(" %s", e.Name)
	if len(e.Value) > 0 {
		str += fmt.Sprintf(" %s", e.Value.ToString())
	} else {
		str += fmt.Sprintf(" %s", e.ExtID.ToString())

		if len(e.NData) > 0 {
			str += fmt.Sprintf(" NDATA %s", e.NData)
		}
	}
	str += ">"
	return str
}

func (n Notation) ToString() string {
	return fmt.Sprintf(`<!NOTATION %s %s>`, n.Name, n.ExtID.ToString())
}

func (p PI) ToString() string {
	str := fmt.Sprintf(`<?%s`, p.Target)
	if len(p.Instruction) > 0 {
		str += fmt.Sprintf(` %s`, p.Instruction)
	}
	str += "?>"
	return str
}

func (c Comment) ToString() string {
	return fmt.Sprintf("<!--%s-->", string(c))
}

func (a AttDef) ToString() string {
	return fmt.Sprintf(" %s %s %s", a.Name, a.Type.ToString(), a.Decl.ToString())
}

func (a AttToken) ToString() string {
	switch a {
	case AttTokenCDATA:
		return "CDATA"
	case AttTokenID:
		return "ID"
	case AttTokenIDREF:
		return "IDREF"
	case AttTokenIDREFS:
		return "IDREFS"
	case AttTokenENTITY:
		return "ENTITY"
	case AttTokenENTITIES:
		return "ENTITIES"
	case AttTokenNMTOKEN:
		return "NMTOKEN"
	case AttTokenNMTOKENS:
		return "NMTOKENS"
	default:
		return ""
	}
}

func (n NotationType) ToString() string {
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

func (e Enum) ToString() string {
	return "(" + strings.Join(e.Cases, "|") + ")"
}

func (d DefaultDeclType) ToString() string {
	switch d {
	case DefaultDeclTypeRequired:
		return "#REQUIRED"
	case DefaultDeclTypeImplied:
		return "#IMPLIED"
	case DefaultDeclTypeFixed:
		return "#FIXED"
	default:
		return ""
	}
}

func (d DefaultDecl) ToString() string {
	str := d.Type.ToString()
	if len(d.AttValue) > 0 {
		str += fmt.Sprintf(" %s", d.AttValue.ToString())
	}
	return str
}

func (EMPTY) ToString() string {
	return "EMPTY"
}

func (ANY) ToString() string {
	return "ANY"
}

func (m Mixed) ToString() string {
	str := `(#PCDATA`
	for _, n := range m.Names {
		str += "|" + n
	}
	str += ")"
	return str
}

func (c Children) ToString() string {
	str := c.ChoiceSeq.ToString()
	if c.Suffix != nil {
		str += string(*c.Suffix)
	}
	return str
}

func (c CP) ToString() string {
	var str string
	if c.ChoiceSeq != nil {
		str = c.ChoiceSeq.ToString()
	} else {
		str = c.Name
	}
	if c.Suffix != nil {
		str += string(*c.Suffix)
	}
	return str
}

func (c Choice) ToString() string {
	str := "("
	for i, cp := range c.CPs {
		if i > 0 {
			str += "|"
		}
		str += cp.ToString()
	}
	str += ")"
	return str
}

func (s Seq) ToString() string {
	str := "("
	for i, cp := range s.CPs {
		if i > 0 {
			str += ","
		}
		str += cp.ToString()
	}
	str += ")"
	return str
}

func (e EntityValue) ToString() string {
	str := `"`
	for _, v := range e {
		if t, ok := v.(Terminal); ok {
			str += t.ToString()
		} else {
			str += fmt.Sprint(v)
		}
	}
	str += `"`
	return str
}

func (a AttValue) ToString() string {
	str := `"`
	for _, v := range a {
		if t, ok := v.(Terminal); ok {
			str += t.ToString()
		} else {
			str += fmt.Sprint(v)
		}
	}
	str += `"`
	return str
}

func (e CharRef) ToString() string {
	return fmt.Sprintf("%s%s;", e.Prefix, e.Value)
}

func (e EntityRef) ToString() string {
	return fmt.Sprintf("&%s;", e.Name)
}

func (e PERef) ToString() string {
	return "%" + fmt.Sprintf(`%s;`, e.Name)
}

func (a Attribute) ToString() string {
	return fmt.Sprintf("%s=%s", a.Name, a.AttValue.ToString())
}

func (a Attributes) ToString() string {
	attrs := make([]string, len(a))
	for i, v := range a {
		attrs[i] = v.ToString()
	}
	return strings.Join(attrs, " ")
}

func (e CData) ToString() string {
	return fmt.Sprintf(`<![CDATA[%s]]>`, e)
}
