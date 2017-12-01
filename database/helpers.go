package database

// QuoteAndEscapeValue wraps escaped double quotes around single composite record values, and
// escapes specific special runes within the value.
// Composite record arrays line format looks like: {"(\"val1\",\"val2\")""}
// Escape code mostly copied from https://github.com/lib/pq/blob/master/encode.go#L160-L196
func QuoteAndEscapeRecordValue(text string) string {
	var c byte

	// byte array to store our final escaped value
	result := make([]byte, 0)
	// start the value with escaped double quotes
	result = append(result, '\\', '"')

	// check each rune for special characters and prepend the proper escape sequence
	for i := 0; i < len(text); i++ {
		c = text[i]

		switch c {
		case '\\':
			result = append(result, '\\', '\\', '\\', '\\')
		case '\n':
			result = append(result, '\\', 'n')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		case '"':
			result = append(result, '\\', '\\', '\\', '"')
		default:
			result = append(result, c)
		}
	}

	// end the value with escaped double quotes
	return string(append(result, '\\', '"'))
}
