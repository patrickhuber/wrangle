package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/templates"
)

type fixtureVariableTokenizer struct {
	tokens []templates.Token
	index  int
}

func (t *fixtureVariableTokenizer) Peek() *templates.Token {
	if t.index == len(t.tokens) {
		return nil
	}
	return &t.tokens[t.index]
}

func (t *fixtureVariableTokenizer) Next() *templates.Token {
	if t.index == len(t.tokens) {
		return nil
	}
	token := &t.tokens[t.index]
	t.index++
	return token
}

func newFixtureTokenizer(tokens []templates.Token) templates.VariableTokenizer {
	return &fixtureVariableTokenizer{
		tokens: tokens,
		index:  0,
	}
}

var _ = Describe("VariableParser", func() {
	It("parses just text", func() {
		tokens := []templates.Token{
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "test",
			},
		}
		tokenizer := newFixtureTokenizer(tokens)

		parser := templates.NewVariableParser()
		output := parser.Parse(tokenizer)
		Expect(output).ToNot(BeNil())
		Expect(len(output.Children)).To(Equal(0))
		Expect(output.Leaf).ToNot(BeNil())
		Expect(output.Leaf.Capture).To(Equal("test"))
	})

	It("parses variable and text", func() {
		tokens := []templates.Token{
			templates.Token{
				TokenType: templates.VariableAstOpen,
				Capture:   "((",
			},
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "test",
			},
			templates.Token{
				TokenType: templates.VariableAstClose,
				Capture:   "))",
			},
		}
		tokenizer := newFixtureTokenizer(tokens)
		parser := templates.NewVariableParser()
		output := parser.Parse(tokenizer)

		Expect(output).ToNot(BeNil())
		Expect(len(output.Children)).To(Equal(1))

		Expect(len(output.Children[0].Children)).To(Equal(3))
		Expect(output.Children[0].Children[0].Leaf).ToNot(BeNil())
		Expect(output.Children[0].Children[0].Leaf.TokenType).To(Equal(templates.VariableAstOpen))

		Expect(output.Children[0].Children[1].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[0].Children[1].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[0].Children[1].Leaf.Capture).To(Equal("test"))

		Expect(output.Children[0].Children[2].Leaf).ToNot(BeNil())
		Expect(output.Children[0].Children[2].Leaf.TokenType).To(Equal(templates.VariableAstClose))
	})

	It("parses nested structures", func() {
		tokens := []templates.Token{
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "before",
			},
			templates.Token{
				TokenType: templates.VariableAstOpen,
				Capture:   "((",
			},
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "before-nest",
			},
			templates.Token{
				TokenType: templates.VariableAstOpen,
				Capture:   "((",
			},
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "nest",
			},
			templates.Token{
				TokenType: templates.VariableAstClose,
				Capture:   "))",
			},
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "after-nest",
			},
			templates.Token{
				TokenType: templates.VariableAstClose,
				Capture:   "))",
			},
			templates.Token{
				TokenType: templates.VariableAstText,
				Capture:   "after",
			},
		}
		tokenizer := newFixtureTokenizer(tokens)
		parser := templates.NewVariableParser()
		output := parser.Parse(tokenizer)
		Expect(output).ToNot(BeNil())
		Expect(len(output.Children)).To(Equal(3))

		// before
		Expect(output.Children[0].Leaf).ToNot(BeNil())
		Expect(output.Children[0].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[0].Leaf.Capture).To(Equal("before"))

		// before nest
		Expect(len(output.Children[1].Children)).To(Equal(5))
		Expect(output.Children[1].Children[0].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[0].Leaf.TokenType).To(Equal(templates.VariableAstOpen))

		Expect(output.Children[1].Children[1].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[1].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[1].Children[1].Leaf.Capture).To(Equal("before-nest"))

		// nest
		Expect(len(output.Children[1].Children[2].Children)).To(Equal(3))

		Expect(output.Children[1].Children[2].Children[0].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[2].Children[0].Leaf.TokenType).To(Equal(templates.VariableAstOpen))

		Expect(output.Children[1].Children[2].Children[1].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[2].Children[1].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[1].Children[2].Children[1].Leaf.Capture).To(Equal("nest"))

		Expect(output.Children[1].Children[2].Children[2].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[2].Children[2].Leaf.TokenType).To(Equal(templates.VariableAstClose))

		// after nest
		Expect(output.Children[1].Children[3].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[3].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[1].Children[3].Leaf.Capture).To(Equal("after-nest"))

		Expect(len(output.Children[1].Children)).To(Equal(5))
		Expect(output.Children[1].Children[4].Leaf).ToNot(BeNil())
		Expect(output.Children[1].Children[4].Leaf.TokenType).To(Equal(templates.VariableAstClose))

		// after
		Expect(output.Children[2].Leaf).ToNot(BeNil())
		Expect(output.Children[2].Leaf.TokenType).To(Equal(templates.VariableAstText))
		Expect(output.Children[2].Leaf.Capture).To(Equal("after"))
	})
})
