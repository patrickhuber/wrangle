package lex

import (
	"fmt"
	"unicode/utf8"

	"github.com/patrickhuber/wrangle/internal/dataptr/token"
)

type Lexer struct {
	peekToken *token.Token
	input     string
	offset    int
	position  int
	column    int
	line      int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) Next() (*token.Token, error) {
	if l.peekToken == nil {
		return l.next()
	}
	tok := l.peekToken
	l.peekToken = nil
	return tok, nil
}

func (l *Lexer) Peek() (*token.Token, error) {
	if l.peekToken != nil {
		return l.peekToken, nil
	}
	p, err := l.next()
	if err != nil {
		return nil, err
	}
	l.peekToken = p
	return p, nil
}

func (l *Lexer) next() (*token.Token, error) {
	r, ok := l.peekRune()
	if !ok {
		return l.token(token.EndOfStream)
	}
	switch {

	case r == '/':
		if err := l.expect('/'); err != nil {
			return nil, err
		}
		return l.token(token.Slash)

	case r == '=':
		if err := l.expect('='); err != nil {
			return nil, err
		}
		return l.token(token.Equal)
	case l.isNumber(r):
		if err := l.eatWhile(l.isNumber); err != nil {
			return nil, err
		}
		return l.token(token.Integer)

	case l.isLetter(r):
		if err := l.eatWhile(l.isLetter); err != nil {
			return nil, err
		}
		return l.token(token.Name)
	}
	return nil, fmt.Errorf("unrecognized token at position %d", l.position)
}

func (l *Lexer) isLetter(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func (l *Lexer) isNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) peekRune() (rune, bool) {
	if l.position >= len(l.input) {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.position:])
	return r, true
}

func (l *Lexer) readRune() (rune, bool) {
	if l.position >= len(l.input) {
		return 0, false
	}
	r, width := utf8.DecodeRuneInString(l.input[l.position:])
	l.position += width
	return r, true
}

// token creates a token of the specified type at the current position
// after calling token, the offset will be advanced
func (l *Lexer) token(ty token.Type) (*token.Token, error) {
	tok := &token.Token{
		Type:     ty,
		Position: l.offset,
		Column:   l.column,
		Line:     l.line,
		Capture:  l.input[l.offset:l.position],
	}
	for i := l.offset; i < l.position; i++ {
		ch := l.input[i]
		if ch == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
	}
	l.offset = l.position
	return tok, nil
}

func (l *Lexer) eatWhile(condition func(r rune) bool) error {
	for {
		ok, err := l.eatIf(condition)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

func (l *Lexer) eatIf(condition func(r rune) bool) (bool, error) {
	p, ok := l.peekRune()
	if !ok {
		return false, nil
	}
	if !condition(p) {
		return false, nil
	}
	err := l.assert(condition)
	return err == nil, err
}

func (l *Lexer) eat(ch rune) (bool, error) {
	p, ok := l.peekRune()
	if !ok {
		return false, nil
	}
	if p != ch {
		return false, nil
	}
	err := l.expect(ch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *Lexer) expect(ch rune) error {
	r, ok := l.readRune()
	if !ok {
		return fmt.Errorf("unable to read rune")
	}
	if r != ch {
		return fmt.Errorf("expected '%c' but found '%c'", ch, r)
	}
	return nil
}

func (l *Lexer) assert(condition func(r rune) bool) error {
	r, ok := l.readRune()
	if !ok {
		return fmt.Errorf("unable to read rune")
	}
	if !condition(r) {
		return fmt.Errorf("expected condition to match '%c'", r)
	}
	return nil
}
