package templates

type TokenType int

const (
	VariableAstOpen  TokenType = 0
	VariableAstClose TokenType = 1
	VariableAstText  TokenType = 2
)

type Token struct {
	TokenType TokenType
	Capture   string
}
