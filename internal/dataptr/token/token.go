package token

type Type string

const (
	Integer     Type = "integer"
	Name        Type = "name"
	Equal       Type = "="
	Slash       Type = "/"
	EndOfStream Type = "EOF"
)

type Token struct {
	Type     Type
	Position int
	Column   int
	Line     int
	Capture  string
}
