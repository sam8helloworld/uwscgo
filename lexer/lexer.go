package lexer

import "github.com/sam8helloworld/uwscgo/token"

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
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	}
	l.readChar()
	return tok
}
