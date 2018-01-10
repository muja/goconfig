package goconfig

const utf8BOM = "\357\273\277"

type parser struct {
	bytes  []byte
	linenr uint
	eof    bool
}

// Parse takes given bytes as configuration file (according to gitconfig syntax)
func Parse(bytes []byte) (map[string]string, uint, error) {
	parser := &parser{bytes, 1, false}
	cfg, err := parser.parse()
	return cfg, parser.linenr, err
}

func (cf *parser) parse() (map[string]string, error) {
	bomPtr := 0
	comment := false
	cfg := map[string]string{}
	name := ""
	var err error
	for {
		c := cf.nextChar()
		if bomPtr != -1 && bomPtr < len(utf8BOM) {
			if c == (utf8BOM[bomPtr] & 0377) {
				bomPtr++
				continue
			} else {
				/* Do not tolerate partial BOM. */
				if bomPtr != 0 {
					return cfg, ErrPartialBOM
				}
				bomPtr = -1
			}
		}
		if c == '\n' {
			if cf.eof {
				return cfg, nil
			}
			comment = false
			continue
		}
		if comment || isspace(c) {
			continue
		}
		if c == '#' || c == ';' {
			comment = true
			continue
		}
		if c == '[' {
			name, err = cf.getSectionKey()
			if err != nil {
				return cfg, err
			}
			name += "."
			continue
		}
		if !isalpha(c) {
			return cfg, ErrInvalidKeyChar
		}
		key := name + string(c)
		value, err := cf.getValue(&key)
		if err != nil {
			return cfg, err
		}
		cfg[key] = value
	}
}

func (cf *parser) nextChar() byte {
	if len(cf.bytes) == 0 {
		cf.eof = true
		return byte('\n')
	}
	c := cf.bytes[0]
	if c == '\r' {
		/* DOS like systems */
		if len(cf.bytes) > 1 && cf.bytes[1] == '\n' {
			cf.bytes = cf.bytes[1:]
			c = '\n'
		}
	}
	if c == '\n' {
		cf.linenr++
	}
	if len(cf.bytes) == 0 {
		cf.eof = true
		cf.linenr++
		c = '\n'
	}
	cf.bytes = cf.bytes[1:]
	return c
}

func (cf *parser) getSectionKey() (string, error) {
	name := ""
	for {
		c := cf.nextChar()
		if cf.eof {
			return "", ErrUnexpectedEOF
		}
		if c == ']' {
			return name, nil
		}
		if isspace(c) {
			return cf.getExtendedSectionKey(name, c)
		}
		if !iskeychar(c) && c != '.' {
			return "", ErrInvalidSectionChar
		}
		name += string(lower(c))
	}
}

// config: [BaseSection "ExtendedSection"]
func (cf *parser) getExtendedSectionKey(name string, c byte) (string, error) {
	for {
		if c == '\n' {
			cf.linenr--
			return "", ErrSectionNewLine
		}
		c = cf.nextChar()
		if !isspace(c) {
			break
		}
	}
	if c != '"' {
		return "", ErrMissingStartQuote
	}
	name += "."
	for {
		c = cf.nextChar()
		if c == '\n' {
			cf.linenr--
			return "", ErrSectionNewLine
		}
		if c == '"' {
			break
		}
		if c == '\\' {
			c = cf.nextChar()
			if c == '\n' {
				cf.linenr--
				return "", ErrSectionNewLine
			}
		}
		name += string(c)
	}
	if cf.nextChar() != ']' {
		return "", ErrMissingClosingBracket
	}
	return name, nil
}

func (cf *parser) getValue(name *string) (string, error) {
	var c byte
	var err error
	var value string

	/* Get the full name */
	for {
		c = cf.nextChar()
		if cf.eof {
			break
		}
		if !iskeychar(c) {
			break
		}
		*name += string(lower(c))
	}

	for c == ' ' || c == '\t' {
		c = cf.nextChar()
	}

	if c != '\n' {
		if c != '=' {
			return "", ErrInvalidKeyChar
		}
		value, err = cf.parseValue()
		if err != nil {
			return "", err
		}
	}
	/*
	 * We already consumed the \n, but we need linenr to point to
	 * the line we just parsed during the call to fn to get
	 * accurate line number in error messages.
	 */
	// cf.linenr--
	// ret := fn(name->buf, value, data);
	// if ret >= 0 {
	// 	cf.linenr++
	// }
	return value, err
}

func (cf *parser) parseValue() (string, error) {
	var quote, comment bool
	var space int

	var value string

	// strbuf_reset(&cf->value);
	for {
		c := cf.nextChar()
		if c == '\n' {
			if quote {
				cf.linenr--
				return "", ErrUnfinishedQuote
			}
			return value, nil
		}
		if comment {
			continue
		}
		if isspace(c) && !quote {
			if len(value) > 0 {
				space++
			}
			continue
		}
		if !quote {
			if c == ';' || c == '#' {
				comment = true
				continue
			}
		}
		for space != 0 {
			value += " "
			space--
		}
		if c == '\\' {
			c = cf.nextChar()
			switch c {
			case '\n':
				continue
			case 't':
				c = '\t'
				break
			case 'b':
				c = '\b'
				break
			case 'n':
				c = '\n'
				break
			/* Some characters escape as themselves */
			case '\\':
				break
			case '"':
				break
			/* Reject unknown escape sequences */
			default:
				return "", ErrInvalidEscapeSequence
			}
			value += string(c)
			continue
		}
		if c == '"' {
			quote = !quote
			continue
		}
		value += string(c)
	}
}

func lower(c byte) byte {
	return c | 0x20
}

func isspace(c byte) bool {
	return c == '\t' || c == ' ' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

func iskeychar(c byte) bool {
	return isalnum(c) || c == '-'
}

func isalnum(c byte) bool {
	return isalpha(c) || isnum(c)
}

func isalpha(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
}

func isnum(c byte) bool {
	return c >= '0' && c <= '9'
}
