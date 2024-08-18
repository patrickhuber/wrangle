package lex_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/dataptr/lex"
	"github.com/patrickhuber/wrangle/internal/dataptr/token"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	type test struct {
		name   string
		str    string
		tokens []token.Type
	}
	tests := []test{
		{"name", "name", []token.Type{token.Name}},
		{"name and path", "name/parent/child", []token.Type{
			token.Name,
			token.Slash,
			token.Name,
			token.Slash,
			token.Name,
		}},
		{"condition", "name/key=value", []token.Type{
			token.Name,
			token.Slash,
			token.Name,
			token.Equal,
			token.Name,
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testLexer(t, test.str, test.tokens)
		})
	}
}

func testLexer(t *testing.T, str string, tokens []token.Type) {
	lexer := lex.New(str)
	var i = 0
	for i = 0; i < len(tokens); i++ {
		actual, err := lexer.Next()
		require.NoError(t, err)
		require.NotNil(t, actual)
		expected := tokens[i]
		require.Equal(t, expected, actual.Type, "expected %s but found %s", expected, actual.Type)
	}
	require.Equal(t, len(tokens), i, "expected token count of %d but found %d", len(tokens), i)
}
