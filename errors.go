package goconfig

import "errors"

// ErrInvalidEscapeSequence indicates that the escape character ('\')
// was followed by an invalid character.
var ErrInvalidEscapeSequence = errors.New("unknown escape sequence")

// ErrUnfinishedQuote indicates that a value has an odd number of (unescaped) quotes
var ErrUnfinishedQuote = errors.New("unfinished quote")

// ErrMissingEquals indicates that an equals sign ('=') was expected but not found
var ErrMissingEquals = errors.New("expected '='")

// ErrPartialBOM indicates that the file begins with a partial UTF8-BOM
var ErrPartialBOM = errors.New("partial UTF8-BOM")

// ErrInvalidKeyChar indicates that there was an invalid key character
var ErrInvalidKeyChar = errors.New("invalid key character")

// ErrInvalidSectionChar indicates that there was an invalid character in section
var ErrInvalidSectionChar = errors.New("invalid character in section")

// ErrUnexpectedEOF indicates that there was an unexpected EOF
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// ErrSectionNewLine indicates that there was a newline in section
var ErrSectionNewLine = errors.New("newline in section")

// ErrMissingStartQuote indicates that there was a missing start quote
var ErrMissingStartQuote = errors.New("missing start quote")

// ErrMissingClosingBracket indicates that there was a missing closing bracket in section
var ErrMissingClosingBracket = errors.New("missing closing section bracket")
