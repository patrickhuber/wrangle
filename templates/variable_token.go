package templates

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