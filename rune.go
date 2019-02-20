package xml

// [0-9]
func isNum(r rune) bool {
	return 0x30 <= r && r <= 0x39
}

// [a-zA-Z]
func isAlpha(r rune) bool {
	// A-Z || a-z
	return 0x41 <= r && r <= 0x5A || 0x61 <= r && r <= 0x7A
}

// #x20 | #x9 | #xD | #xA
func isSpace(r rune) bool {
	return r == 0x20 || r == 0x9 || r == 0xD || r == 0xA
}

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
