package main

import (
	xml "github.com/matsune/go-xml"
)

func main() {
	str := `
	<?xml version="1.0" encoding="UTF-8" standalone="yes"?><!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
				"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd" [ <!ELEMENT code (#PCDATA)> <!NOTATION vrml PUBLIC "VRML 1.0"> <!ATTLIST code lang NOTATION (vrml) #REQUIRED>
		%peref;   <!ATTLIST task status (important|normal) #REQUIRED>
<!ATTLIST task status NMTOKEN #FIXED "monthly">  <!ATTLIST description xml:lang NMTOKEN #FIXED "en">
<!--'student' must have   	at least one child element of type 'subject'-->
<!ELEMENT student (subject+)>
	  ]>
	  <root xmlns:h="http://www.w3.org/TR/html4/"
xmlns:f="https://www.w3schools.com/furniture">

	  <h:table xmlns:h="http://www.w3.org/TR/html4/">
  <h:tr>
    <h:td>Apples</h:td>
    <h:td>Bananas</h:td>
  </h:tr>
</h:table></root>


<!--'student' must have   	at least one child element of type 'subject'-->
`
	p := xml.NewParser(str)
	if x, err := p.Parse(); err != nil {
		panic(err)
	} else {
		f := xml.NewFormatter()
		f.Format(x)
	}
}
