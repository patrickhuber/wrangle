package templates

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
)

// Template - contains a go object document that contains placeholders for variables
type Template interface {
	Evaluate(resolvers ...VariableResolver) (interface{}, error)
}

type template struct {
	document interface{}
}

// NewTemplate - Creates a new template with the given document parameter
func NewTemplate(document interface{}) Template {
	return &template{document: document}
}

// Evaluate - Evaluates the tempalte using the variable resolver for variable lookup
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
	case ([]string):
		return evaluateSliceOfString(t, resolver)
	case ([]interface{}):
		return evaluateSliceOfInterface(t, resolver)
	}
	return nil, fmt.Errorf("unable to evaluate template. Invalid type '%v'", reflect.TypeOf(document))
}

func evaluateString(template string, resolver VariableResolver) (interface{}, error) {
	re := regexp.MustCompile(`\(\((?P<key>[^)]*)\)\)`)
	index := 0
	result := bytes.Buffer{}
	for _, submatches := range re.FindAllStringSubmatchIndex(template, -1) {
		matchStart := submatches[0]
		matchEnd := submatches[1]
		subMatchStart := submatches[2]
		subMatchEnd := submatches[3]

		key := template[subMatchStart:subMatchEnd]
		value, err := resolver.Get(string(key))
		if err != nil {
			return nil, err
		}

		// return the value if it is the only match
		if matchStart == 0 && matchEnd == len(template) {
			return value, nil
		}

		// check the value type
		// if it is non-string, return a failure
		v, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("unable to evaluate template. Multiple values detected in source string and key '%s' resulted in a variable lookup of type '%v'. Make sure the source is a single variable ((variable)) or the lookup value is of type string", key, reflect.TypeOf(v))
		}

		// append the string previous to the match to the result
		// up to the start of the match
		result.WriteString(template[index:matchStart])
		// update the index
		index = submatches[1]
		// append the lookup result
		result.WriteString(v)
	}
	if index < len(template) {
		result.WriteString(template[index:])
	}
	return result.String(), nil
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
