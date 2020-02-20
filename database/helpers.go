package database

import (
	"bytes"
	"fmt"
)

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

// ParseCustomType parses a Postgres custom type retrieved from the database.
func ParseCustomType(src, del []byte) (elems [][]byte, err error) {
	var i int

	// the first byte should be ( to indicate the start of a custom type
	if len(src) < 1 || src[0] != '(' {
		return nil, fmt.Errorf("unable to parse custom type; unexpected %q at offset %d", src[i], i)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '(':
			// increment and move on since this char is not part of an element
			i++
		case ')':
			goto Close
		default:
			// we have found an element
			break Open
		}
	}

Element:
	// start iterating over the bytes of the element
	for i < len(src) {
		switch src[i] {
		case '"':
			// we found the start of a string
			var elem = []byte{}
			var escape bool
			// increment i before looking at the parts of the element
			// this will drop the initial " since we don't want it as
			// part of the element
			for i++; i < len(src); i++ {
				if escape {
					// an escape sequence ( "" or \" ) was found
					// append the preceding bytes since it should
					// be part of the element
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						// append this part of the element
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						// this case could be the end of the element
						// or it could be the start of "" incicating that
						// we have found a quote that needs to be part of
						// the element
						if src[i+1] == '"' {
							// not the end of the element since we found ""
							escape = true
						} else {
							// this is the end of the element
							elems = append(elems, elem)
							i++
							break Element
						}
					}
				}
			}
		default:
			// we found the start of a non string element so we can append
			// each byte until we find the delimiter (del param) or the `)`
			// indicating the end of the custom type
			for start := i; i < len(src); i++ {
				// don't increment start so we know what cutset to take
				// from src to use as the part of the element
				if bytes.HasPrefix(src[i:], del) || src[i] == ')' {
					elem := src[start:i]
					if len(elem) == 0 {
						elem = nil
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)

					break Element
				}
			}
		}
	}

	for i < len(src) {
		// if we have a delimiter we need to continue grabbing
		// elements of the custom type
		if bytes.HasPrefix(src[i:], del) {
			i += len(del)
			goto Element
		} else if src[i] == ')' {
			// found the end, increment i so we break out of our looping
			i++
		} else {
			return nil, fmt.Errorf("unable to parse custom type; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	// if we get here that indicates an empty custom type ()
	for i < len(src) {
		if src[i] == ')' {
			i++
		} else {
			return nil, fmt.Errorf("unable to parse custom type; unexpected %q at offset %d", src[i], i)
		}
	}

	return
}
