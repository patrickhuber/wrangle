package parse

import (
	"fmt"
	"strconv"

	"github.com/patrickhuber/wrangle/internal/dataptr/ast"
	"github.com/patrickhuber/wrangle/internal/dataptr/lex"
	"github.com/patrickhuber/wrangle/internal/dataptr/token"
)

func Parse(str string) (ast.DataPointer, error) {
	lexer := lex.New(str)
	return parse(lexer)
}

func parse(lexer *lex.Lexer) (ast.DataPointer, error) {
	var segments []ast.Segment
	for {
		// can be a single segment
		segment, err := parseSegment(lexer)
		if err != nil {
			return ast.DataPointer{}, err
		}
		segments = append(segments, segment)

		// or multiple
		ok, err := eat(lexer, token.Slash)
		if err != nil {
			return ast.DataPointer{}, err
		}
		if !ok {
			break
		}
	}
	return ast.DataPointer{
		Segments: segments,
	}, nil
}

func parseSegment(lexer *lex.Lexer) (ast.Segment, error) {
	tok, err := lexer.Peek()
	if err != nil {
		return nil, err
	}

	// we have an integer
	if tok.Type == token.Integer {
		i, err := parseInteger(lexer)
		if err != nil {
			return nil, err
		}
		return ast.Index{
			Index: i,
		}, nil
	}

	// otherwise this is a name
	name, err := parseName(lexer)
	if err != nil {
		return nil, err
	}

	// if no equal sign, this is an element
	ok, err := eat(lexer, token.Equal)
	if err != nil {
		return nil, err
	}
	if !ok {
		return ast.Element{
			Name: name,
		}, nil
	}

	// this is a constraint 'key=value'
	value, err := parseName(lexer)
	if err != nil {
		return nil, err
	}
	return ast.Constraint{
		Key:   name,
		Value: value,
	}, nil
}

func parseInteger(lexer *lex.Lexer) (int, error) {
	tok, err := expect(lexer, token.Integer)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(tok.Capture, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func parseName(lexer *lex.Lexer) (string, error) {
	tok, err := expect(lexer, token.Name)
	if err != nil {
		return "", err
	}
	return tok.Capture, nil
}

func eat(lexer *lex.Lexer, ty token.Type) (bool, error) {
	p, err := lexer.Peek()
	if err != nil {
		return false, err
	}
	if p.Type != ty {
		return false, nil
	}
	_, err = lexer.Next()
	if err != nil {
		return false, err
	}
	return true, nil
}

func expect(lexer *lex.Lexer, ty token.Type) (*token.Token, error) {
	tok, err := lexer.Next()
	if err != nil {
		return nil, err
	}
	if tok.Type != ty {
		return nil, parseError(tok, token.Name)
	}
	return tok, nil
}

func parseError(tok *token.Token, expected token.Type) error {
	return fmt.Errorf("error parsing at line: %d column: %d. Expected: %s found: %s", tok.Line, tok.Column, expected, tok.Type)
}
