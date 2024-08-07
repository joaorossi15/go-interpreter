package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input   string
	pos     int  // current position in input, pointing to current char
	readPos int  // current read position in input, pointing after current char
	ch      byte // current char examinated
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // initialize lexer
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos += 1
}

func (l *Lexer) NextToken() token.Token {
	var tk token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tk = token.Token{Type: token.EQ, Literal: string(l.ch) + string(l.ch)}
		} else {
			tk = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tk = newToken(token.PLUS, l.ch)
	case '(':
		tk = newToken(token.LPAREN, l.ch)
	case ')':
		tk = newToken(token.RPAREN, l.ch)
	case '{':
		tk = newToken(token.LBRACE, l.ch)
	case '}':
		tk = newToken(token.RBRACE, l.ch)
	case ',':
		tk = newToken(token.COMMA, l.ch)
	case ';':
		tk = newToken(token.SEMICOLON, l.ch)
	case '*':
		tk = newToken(token.ASTERISK, l.ch)
	case '/':
		tk = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tk = token.Token{Type: token.NOT_EQ, Literal: string(l.ch) + string('=')}
		} else {
			tk = newToken(token.BANG, l.ch)
		}
	case '-':
		tk = newToken(token.MINUS, l.ch)
	case '<':
		tk = newToken(token.LT, l.ch)
	case '>':
		tk = newToken(token.GT, l.ch)
	case 0:
		tk.Literal = ""
		tk.Type = "EOF"
	default:
		if isLetter(l.ch) {
			tk.Literal = l.readIdentifier()
			tk.Type = token.LookupIdent(tk.Literal)
			return tk
		} else if isDigit(l.ch) {
			tk.Literal = l.readNumber()
			tk.Type = token.INT
			return tk
		} else {
			tk = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tk
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	pos := l.pos

	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}
