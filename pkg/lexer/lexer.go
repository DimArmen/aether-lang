package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType string

const (
	// Special tokens
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"

	// Identifiers and literals
	TokenIdent  TokenType = "IDENT"
	TokenString TokenType = "STRING"
	TokenNumber TokenType = "NUMBER"
	TokenBool   TokenType = "BOOL"

	// Keywords
	TokenResource      TokenType = "RESOURCE"
	TokenVariable      TokenType = "VARIABLE"
	TokenOutput        TokenType = "OUTPUT"
	TokenModule        TokenType = "MODULE"
	TokenAgent         TokenType = "AGENT"
	TokenProvider      TokenType = "PROVIDER"
	TokenScript        TokenType = "SCRIPT"
	TokenFor           TokenType = "FOR"
	TokenIn            TokenType = "IN"
	TokenIf            TokenType = "IF"
	TokenElse          TokenType = "ELSE"
	TokenMatch         TokenType = "MATCH"
	TokenLet           TokenType = "LET"
	TokenFunction      TokenType = "FUNCTION"
	TokenReturn        TokenType = "RETURN"
	TokenImport        TokenType = "IMPORT"
	TokenExport        TokenType = "EXPORT"
	TokenTest          TokenType = "TEST"
	TokenAssert        TokenType = "ASSERT"
	TokenDependsOn     TokenType = "DEPENDS_ON"
	TokenCount         TokenType = "COUNT"
	TokenForEach       TokenType = "FOREACH"
	TokenTypeKeyword   TokenType = "TYPE"
	TokenSecret        TokenType = "SECRET"
	TokenBackend       TokenType = "BACKEND"
	TokenTrue          TokenType = "TRUE"
	TokenFalse         TokenType = "FALSE"
	TokenNull          TokenType = "NULL"
	TokenAnd           TokenType = "AND"
	TokenOr            TokenType = "OR"
	TokenNot           TokenType = "NOT"

	// Operators
	TokenAssign   TokenType = "="
	TokenPlus     TokenType = "+"
	TokenMinus    TokenType = "-"
	TokenStar     TokenType = "*"
	TokenSlash    TokenType = "/"
	TokenPercent  TokenType = "%"
	TokenPower    TokenType = "**"
	TokenEq       TokenType = "=="
	TokenNotEq    TokenType = "!="
	TokenLt       TokenType = "<"
	TokenGt       TokenType = ">"
	TokenLtEq     TokenType = "<="
	TokenGtEq     TokenType = ">="
	TokenArrow    TokenType = "=>"
	TokenPipe     TokenType = "|>"
	TokenQuestion TokenType = "?"
	TokenColon    TokenType = ":"
	TokenDot      TokenType = "."
	TokenComma    TokenType = ","
	TokenSemi     TokenType = ";"

	// Delimiters
	TokenLParen   TokenType = "("
	TokenRParen   TokenType = ")"
	TokenLBrace   TokenType = "{"
	TokenRBrace   TokenType = "}"
	TokenLBracket TokenType = "["
	TokenRBracket TokenType = "]"
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer tokenizes Aether source code
type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next reading position
	ch           byte // current char
	line         int
	column       int
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()
	l.skipComments()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenEq, Literal: "==", Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TokenArrow, Literal: "=>", Line: l.line, Column: l.column - 1}
		} else {
			tok = Token{Type: TokenAssign, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			tok = Token{Type: TokenPower, Literal: "**", Line: l.line, Column: l.column - 1}
		} else {
			tok = Token{Type: TokenStar, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '/':
		tok = Token{Type: TokenSlash, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '%':
		tok = Token{Type: TokenPercent, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenNotEq, Literal: "!=", Line: l.line, Column: l.column - 1}
		} else {
			tok = Token{Type: TokenNot, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenLtEq, Literal: "<=", Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '<' {
			// Handle heredoc <<-EOF
			return l.readHeredoc()
		} else {
			tok = Token{Type: TokenLt, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenGtEq, Literal: ">=", Line: l.line, Column: l.column - 1}
		} else {
			tok = Token{Type: TokenGt, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '|':
		if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TokenPipe, Literal: "|>", Line: l.line, Column: l.column - 1}
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '?':
		tok = Token{Type: TokenQuestion, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ':':
		tok = Token{Type: TokenColon, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '.':
		tok = Token{Type: TokenDot, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ',':
		tok = Token{Type: TokenComma, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ';':
		tok = Token{Type: TokenSemi, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '(':
		tok = Token{Type: TokenLParen, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ')':
		tok = Token{Type: TokenRParen, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '{':
		tok = Token{Type: TokenLBrace, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '}':
		tok = Token{Type: TokenRBrace, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '[':
		tok = Token{Type: TokenLBracket, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ']':
		tok = Token{Type: TokenRBracket, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	for {
		if l.ch == '/' && l.peekChar() == '/' {
			// Single-line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			l.skipWhitespace()
		} else if l.ch == '/' && l.peekChar() == '*' {
			// Multi-line comment
			l.readChar() // skip /
			l.readChar() // skip *
			for {
				if l.ch == 0 {
					break
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar() // skip *
					l.readChar() // skip /
					break
				}
				l.readChar()
			}
			l.skipWhitespace()
		} else {
			break
		}
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// Handle decimals
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	// Handle scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var result strings.Builder
	l.readChar() // skip opening "

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '"':
				result.WriteByte('"')
			case '\\':
				result.WriteByte('\\')
			default:
				result.WriteByte(l.ch)
			}
			l.readChar()
		} else if l.ch == '$' && l.peekChar() == '{' {
			// String interpolation - for now just include as is
			// TODO: Handle interpolation properly
			result.WriteByte(l.ch)
			l.readChar()
		} else {
			result.WriteByte(l.ch)
			l.readChar()
		}
	}

	return result.String()
}

func (l *Lexer) readHeredoc() Token {
	tok := Token{Type: TokenString, Line: l.line, Column: l.column}
	l.readChar() // skip first <
	l.readChar() // skip second <
	if l.ch == '-' {
		l.readChar() // skip -
	}

	// Read delimiter
	var delimiter strings.Builder
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		delimiter.WriteByte(l.ch)
		l.readChar()
	}

	// Skip to next line
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.readChar() // skip newline

	// Read until delimiter
	var content strings.Builder
	delim := delimiter.String()
	for {
		lineStart := l.position
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		line := l.input[lineStart:l.position]

		if strings.TrimSpace(line) == delim {
			break
		}

		content.WriteString(line)
		content.WriteByte('\n')

		if l.ch == 0 {
			break
		}
		l.readChar() // skip newline
	}

	tok.Literal = content.String()
	return tok
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

var keywords = map[string]TokenType{
	"resource":   TokenResource,
	"variable":   TokenVariable,
	"output":     TokenOutput,
	"module":     TokenModule,
	"agent":      TokenAgent,
	"provider":   TokenProvider,
	"script":     TokenScript,
	"for":        TokenFor,
	"in":         TokenIn,
	"if":         TokenIf,
	"else":       TokenElse,
	"match":      TokenMatch,
	"let":        TokenLet,
	"function":   TokenFunction,
	"return":     TokenReturn,
	"import":     TokenImport,
	"export":     TokenExport,
	"test":       TokenTest,
	"assert":     TokenAssert,
	"depends_on": TokenDependsOn,
	"count":      TokenCount,
	"for_each":   TokenForEach,
	"type":       TokenTypeKeyword,
	"secret":     TokenSecret,
	"backend":    TokenBackend,
	"true":       TokenTrue,
	"false":      TokenFalse,
	"null":       TokenNull,
	"and":        TokenAnd,
	"or":         TokenOr,
	"not":        TokenNot,
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenIdent
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}",
		t.Type, t.Literal, t.Line, t.Column)
}
