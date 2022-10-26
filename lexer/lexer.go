package lexer

import (
	"github.com/sam8helloworld/uwscgo/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok = token.Token{}

	l.skipWhiteSpace()
	switch l.ch {
	case '=':
		tok = token.Token{
			Type:    token.ASSIGN,
			Literal: string(l.ch),
		}
	case '+':
		tok = token.Token{
			Type:    token.PLUS,
			Literal: string(l.ch),
		}
	case '-':
		tok = token.Token{
			Type:    token.MINUS,
			Literal: string(l.ch),
		}
	case '*':
		tok = token.Token{
			Type:    token.ASTERISK,
			Literal: string(l.ch),
		}
	case '/':
		tok = token.Token{
			Type:    token.SLASH,
			Literal: string(l.ch),
		}
	case '(':
		tok = token.Token{
			Type:    token.LEFT_PARENTHESIS,
			Literal: string(l.ch),
		}
	case ')':
		tok = token.Token{
			Type:    token.RIGHT_PARENTHESIS,
			Literal: string(l.ch),
		}
	case '{':
		tok = token.Token{
			Type:    token.LEFT_BRACKET,
			Literal: string(l.ch),
		}
	case '}':
		tok = token.Token{
			Type:    token.RIGHT_BRACKET,
			Literal: string(l.ch),
		}
	case ',':
		tok = token.Token{
			Type:    token.COMMA,
			Literal: string(l.ch),
		}
	case '\r':
	case '\n':
		tok = token.Token{
			Type:    token.EOL,
			Literal: string(l.ch),
		}
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.Token{
				Type:    token.ILLEGAL,
				Literal: string(l.ch),
			}
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
