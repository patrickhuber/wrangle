package templates

import (
	"strings"

	"github.com/patrickhuber/wrangle/collections"
)

type variableTokenizer struct {
	position int
	state    int
	input    string
	capture  strings.Builder
	emit     *Token
	queue    collections.Queue
}

type VariableTokenizer interface {
	Next() *Token
	Peek() *Token
}

// NewVariableTokenizer creates a new variable tokenizer over the given string
func NewVariableTokenizer(input string) VariableTokenizer {
	return &variableTokenizer{
		position: 0,
		input:    input,
		state:    0,
		queue:    collections.NewQueue(),
	}
}

func (t *variableTokenizer) Next() *Token {
	if !t.queue.Empty() {
		return t.queue.Dequeue().(*Token)
	}
	token := t.Peek()
	if token != nil {
		t.queue.Dequeue()
	}
	return token
}

func (t *variableTokenizer) Peek() *Token {
	if !t.queue.Empty() {
		return t.queue.Peek().(*Token)
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

			if t.capture.Len() > 0 {
				t.queue.Enqueue(&Token{
					TokenType: Text,
					Capture:   t.capture.String(),
				})
				t.capture.Reset()
			}

			t.queue.Enqueue(
				&Token{
					TokenType: OpenVariable,
					Capture:   "((",
				})

			return t.queue.Peek().(*Token)

		case 2:
			t.state = 0

			if ch != ')' {
				t.capture.WriteRune(')')
				t.capture.WriteRune(ch)
				break
			}

			if t.capture.Len() > 0 {
				t.queue.Enqueue(&Token{
					TokenType: Text,
					Capture:   t.capture.String(),
				})
				t.capture.Reset()
			}

			t.queue.Enqueue(
				&Token{
					TokenType: CloseVariable,
					Capture:   "))",
				})

			return t.queue.Peek().(*Token)
		}
	}

	if t.capture.Len() == 0 {
		return nil
	}

	t.queue.Enqueue(
		&Token{
			TokenType: Text,
			Capture:   t.capture.String(),
		})
	t.capture.Reset()
	return t.queue.Peek().(*Token)
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
