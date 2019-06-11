package templates

import (		
	"strings"
	"fmt"
	"reflect"	
)

// Template - contains a go object document that contains placeholders for variables
type Template interface {
	Evaluate(resolvers ...VariableResolver) (interface{}, error)
}

type template struct {
	document interface{}
	macroManager MacroManager
}

// NewTemplate - Creates a new template with the given document parameter
func NewTemplate(document interface{}, macroManager MacroManager) Template {
	return &template{
		document: document,
		macroManager: macroManager}
}

// Evaluate - Evaluates the template using the variable resolver for variable lookup
func (template template) Evaluate(resolvers ...VariableResolver) (interface{}, error) {
	var document = template.document
	var err error
	for _, resolver := range resolvers {
		document, err = evaluate(document, resolver)
		if err != nil {
			return nil, err
		}
	}
	return document, nil
}

// evaluate - dispatches the template evaluation by type to a type specific resolver
// this method is also called recursively by map and array type evaluators
func evaluate(document interface{}, resolver VariableResolver) (interface{}, error) {
	switch t := document.(type) {
	case (string):
		return evaluateString(t, resolver)
	case (map[string]string):
		return evaluateMapStringOfString(t, resolver)
	case (map[string]interface{}):
		return evaluateMapStringOfInterface(t, resolver)
	case (map[interface{}]interface{}):
		return evaluateMapInterfaceOfInterface(t, resolver)
	case ([]string):
		return evaluateSliceOfString(t, resolver)
	case ([]interface{}):
		return evaluateSliceOfInterface(t, resolver)
	}
	return nil, fmt.Errorf("unable to evaluate template. Invalid type '%v'", reflect.TypeOf(document))
}

func evaluateString(template string, resolver VariableResolver) (interface{}, error) {
	tokenizer := NewVariableTokenizer(template)
	parser := NewVariableParser()
	ast := parser.Parse(tokenizer)
	
	return evaluateVariableAst(ast, resolver)
}

func evaluateVariableAst(ast *VariableAst, resolver VariableResolver)(interface{}, error){
	// leaf text node
	if ast.Leaf != nil && ast.Leaf.TokenType == VariableAstText {
		return ast.Leaf.Capture, nil
	}

	if len(ast.Children ) == 0 {
		return nil, fmt.Errorf("invalid ast detected, no children of non text node is invalid")
	}
		
	value := ""
	isClosure := false
	for i, n := range ast.Children{

		if n.Leaf != nil{
			if n.Leaf.TokenType == VariableAstOpen {
				if i == 0{
					isClosure = true
				}
				continue
			}
			if n.Leaf.TokenType == VariableAstClose {
				continue
			}
		}

		vTemp, err := evaluateVariableAst(n, resolver)
		if err != nil {
			return nil, err
		}
		v, ok := vTemp.(string)
		if !ok {
			if len(value) > 0 {
				return nil, fmt.Errorf("resolved value is a structure. values of mixed strings and structures are not allowed")
			}
			return vTemp, nil
		}
		value += v
	}
	if !isClosure {
		return value, nil
	}
	if !strings.HasPrefix(value, "/"){
		value = "/" + value
	}
	return resolver.Get(value)
}

func evaluateMapStringOfString(template map[string]string, resolver VariableResolver) (interface{}, error) {
	transformMap := make(map[string]interface{})
	for k, v := range template {
		value, err := evaluateString(v, resolver)
		if err != nil {
			return nil, err
		}
		transformMap[k] = value
	}
	return transformMap, nil
}

func evaluateMapStringOfInterface(template map[string]interface{}, resolver VariableResolver) (interface{}, error) {
	for k, v := range template {
		value, err := evaluate(v, resolver)
		if err != nil {
			return nil, err
		}
		template[k] = value
	}
	return template, nil
}

func evaluateMapInterfaceOfInterface(template map[interface{}]interface{}, resolver VariableResolver) (interface{}, error) {
	for k, v := range template {
		value, err := evaluate(v, resolver)
		if err != nil {
			return nil, err
		}
		template[k] = value
	}
	return template, nil
}

func evaluateSliceOfString(template []string, resolver VariableResolver) (interface{}, error) {
	// make(type, len, capacity)
	transformSlice := make([]interface{}, 0, len(template))
	for _, v := range template {
		value, err := evaluateString(v, resolver)
		if err != nil {
			return nil, err
		}
		transformSlice = append(transformSlice, value)
	}
	return transformSlice, nil
}

func evaluateSliceOfInterface(template []interface{}, resolver VariableResolver) (interface{}, error) {
	for i, v := range template {
		value, err := evaluate(v, resolver)
		if err != nil {
			return nil, err
		}
		template[i] = value
	}
	return template, nil
}
