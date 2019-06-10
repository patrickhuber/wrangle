package templates

type VariableAst struct {
	Children []*VariableAst
	Leaf     *Token
}
