package baisl

import (
	"os"
	"strings"
)

// Represents a source file that is being lexed
type SourceFile struct {
	path    string
	content []byte
	len     int

	// Current position in the file at lexing time
	index  int
	line   int
	column int
}

func GetSourceFile(path string) (SourceFile, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return SourceFile{}, err
	}

	len := len(file)
	return SourceFile{
		path:    path,
		content: file,
		len:     len,
		index:   0,
		line:    1,
		column:  0,
	}, nil
}

const spaceChars = " \t\n\r\f\v"

func isSpace(c byte) bool {
	return strings.Contains(spaceChars, string(c))
}

const newLineChars = "\n\r"

func isNewLine(c byte) bool {
	return strings.Contains(newLineChars, string(c))
}

// Returns the next character without incrementing the position
func (file *SourceFile) PeekNextChar() (byte, bool) {
	if file.index >= file.len {
		return 0, false
	}

	return file.content[file.index], true
}

// Returns the next character and increments the position
func (file *SourceFile) EatNextChar() (byte, bool) {
	c, ok := file.PeekNextChar()
	if !ok {
		return 0, false
	}

	if isNewLine(c) {
		file.line++
		file.column = 0
	} else {
		file.column++
	}

	file.index++

	return c, true
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

// Returns the next token in the source file
func (file *SourceFile) GetNextToken() Token {
	next, ok := file.EatNextChar()
	if !ok {
		return Token{
			TType: TokenType_EOF,
			Location: SourceLocation{
				Path:   file.path,
				Line:   file.line,
				Column: file.column,
			},
			HasValue: false,
		}
	}
	for isSpace(next) {
		next, ok = file.EatNextChar()
		if !ok {
			return Token{
				TType: TokenType_EOF,
				Location: SourceLocation{
					Path:   file.path,
					Line:   file.line,
					Column: file.column,
				},
				HasValue: false,
			}
		}
	}

	startLoc := SourceLocation{
		Path:   file.path,
		Line:   file.line,
		Column: file.column,
	}

	// Single line token types, split into concrete branches for optimization (probably premature)
	if next == '{' {
		return Token{
			TType:    TokenType_LBRACE,
			Location: startLoc,
			HasValue: false,
		}
	}

	if next == '}' {
		return Token{
			TType:    TokenType_RBRACE,
			Location: startLoc,
			HasValue: false,
		}
	}

	if next == '(' {
		return Token{
			TType:    TokenType_LPAREN,
			Location: startLoc,
			HasValue: false,
		}
	}

	if next == ')' {
		return Token{
			TType:    TokenType_RPAREN,
			Location: startLoc,
			HasValue: false,
		}
	}

	if next == ':' {
		return Token{
			TType:    TokenType_COLON,
			Location: startLoc,
			HasValue: false,
		}
	}

	if next == '/' {
		nextNext, ok := file.PeekNextChar()
		if ok && nextNext == '/' {
			nextNext, ok = file.EatNextChar()
			for !isNewLine(nextNext) && ok {
				file.EatNextChar()
			}

			return file.GetNextToken()
		}
	}

	if IsAlpha(next) {
		value := string(next)
		next, ok = file.PeekNextChar()
		for isAlphaNumeric(next) && ok {
			_, _ = file.EatNextChar()
			value += string(next)
			next, ok = file.PeekNextChar()
		}

		tokenType, exists := KeywordToTokenType[value]
		if exists {
			return Token{
				TType:    tokenType,
				Location: startLoc,
				HasValue: false,
			}
		}

		return Token{
			TType:    TokenType_IDENTIFIER,
			Location: startLoc,
			HasValue: true,
			Value:    value,
		}
	}

	if IsNumeric(next) {
		value := string(next)
		next, ok = file.PeekNextChar()
		for IsNumeric(next) && ok {
			_, _ = file.EatNextChar()
			value += string(next)
			next, ok = file.PeekNextChar()
		}

		return Token{
			TType:    TokenType_NUMBER,
			Location: startLoc,
			HasValue: true,
			Value:    value,
		}
	}

	return Token{
		TType:    TokenType_UNKNOWN,
		Location: startLoc,
		HasValue: true,
		Value:    string(next),
	}
}
