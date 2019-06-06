package templates

import "strings"

type tokenizer struct {
	position int
	state    int
	input    string
	capture  strings.Builder
	emit     *Token
}

type Tokenizer interface {
	Next() *Token
}

func (t *tokenizer) Next() *Token {

	if t.emit != nil {
		emit := t.emit
		t.emit = nil
		return emit
	}

	if t.position == len([]rune(t.input)) {
		return nil
	}

	start := t.position
	for _, ch := range t.input[start:] {
		t.position++
		switch t.state {
		case 0:
			if ch == '(' {
				t.state = 1
			} else if ch == ')' {
				t.state = 2
			} else {
				t.capture.WriteRune(ch)
			}
			break
		case 1:
			t.state = 0

			if ch != '(' {
				t.capture.WriteRune('(')
				t.capture.WriteRune(ch)
				break
			}

			if t.capture.Len() == 0 {
				return &Token{
					TokenType: OpenVariable,
					Capture:   "((",
				}
			}
			t.emit = &Token{
				TokenType: OpenVariable,
				Capture:   "((",
			}

			capture := t.capture.String()
			t.capture.Reset()
			return &Token{
				TokenType: Text,
				Capture:   capture,
			}

		case 2:
			t.state = 0

			if ch != ')' {
				t.capture.WriteRune(')')
				t.capture.WriteRune(ch)
				break
			}

			if t.capture.Len() == 0 {
				return &Token{
					TokenType: CloseVariable,
					Capture:   "))",
				}
			}

			t.emit = &Token{
				TokenType: CloseVariable,
				Capture:   "))",
			}

			capture := t.capture.String()
			t.capture.Reset()
			return &Token{
				TokenType: Text,
				Capture:   capture,
			}
		}
	}

	if t.capture.Len() == 0 {
		return nil
	}

	capture := t.capture.String()
	t.capture.Reset()
	return &Token{
		TokenType: Text,
		Capture:   capture,
	}
}

func NewTokenizer(input string) Tokenizer {
	return &tokenizer{
		position: 0,
		input:    input,
		state:    0,
	}
}

type TokenType int

const (
	OpenVariable  TokenType = 0
	CloseVariable TokenType = 1
	Text          TokenType = 2
)

type Token struct {
	TokenType TokenType
	Capture   string
}
