package templates 

import (
	"github.com/patrickhuber/wrangle/collections"
)

type variableParser struct{

}

type VariableParser interface{
	Parse(tokenizer VariableTokenizer) *VariableAst
}

type VariableAst struct{
	Children []*VariableAst
	Leaf *Token
}

func NewVariableParser() VariableParser{
	return &variableParser{

	}
}

func (p *variableParser) Parse(tokenizer VariableTokenizer) *VariableAst{
	stack := collections.NewStack()
	ast := &VariableAst{
		Children: []*VariableAst{},
	}
	for {
		token := tokenizer.Next()
		
		if token == nil{
			break
		}

		switch token.TokenType{
		case Text:
			peek := tokenizer.Peek()
			if peek == nil && len(ast.Children) == 0{
				ast.Leaf = token
				break
			}
			ast.Children = append(ast.Children, &VariableAst{Leaf: token})
			break

		case OpenVariable:	
			ast.Children = append(ast.Children, &VariableAst{ Leaf: token})
			stack.Push(ast)
			ast = &VariableAst{
				Children: []*VariableAst{},
			}
			break

		case CloseVariable:				
			oldAst := ast		
			ast  = stack.Pop().(*VariableAst)
			ast.Children = append(ast.Children, oldAst)
			ast.Children = append(ast.Children, &VariableAst{ Leaf: token})
			break
		}
	}

	return ast
}