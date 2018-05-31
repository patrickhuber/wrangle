package templates

import (
	"bytes"
	"regexp"
)

// Template - contains a go object document that contains placeholders for variables
type Template interface {
	Evaluate(resolver VariableResolver) interface{}
}

type template struct {
	document interface{}
}

// NewTemplate - Creates a new template with the given document parameter
func NewTemplate(document interface{}) Template {
	return &template{document: document}
}

// Evaluate - Evaluates the tempalte using the variable resolver for variable lookup
func (template template) Evaluate(resolver VariableResolver) interface{} {
	return evaluate(template.document, resolver)
}

// evaluate - dispatches the template evaluation by type to a type specific resolver
// this method is also called recursively by map and array type evaluators
func evaluate(document interface{}, resolver VariableResolver) interface{} {
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
	return document
}

func evaluateString(template string, resolver VariableResolver) interface{} {
	re := regexp.MustCompile(`\(\((?P<key>[^)]*)\)\)`)
	index := 0
	result := bytes.Buffer{}
	for _, submatches := range re.FindAllStringSubmatchIndex(template, -1) {
		matchStart := submatches[0]
		matchEnd := submatches[1]
		subMatchStart := submatches[2]
		subMatchEnd := submatches[3]

		key := template[subMatchStart:subMatchEnd]
		value := resolver.Get(string(key))

		// return the value if it is the only match
		if matchStart == 0 && matchEnd == len(template) {
			return value
		}

		// check the value type
		// if it is non-string, return the template
		v, ok := value.(string)
		if !ok {
			return template
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
	return result.String()
}

func evaluateString1(template string, resolver VariableResolver) interface{} {
	// create a regex that finds the pattern ((key))
	// where key is the path to the variable in the variable store
	re := regexp.MustCompile(`\(\((?P<key>[^)]*)\)\)`)

	// replace all strings will find the pattern in the string and call the match function
	// for each occurance
	return re.ReplaceAllStringFunc(template, func(value string) string {

		result := []byte{}

		// it will expand just the capture group essintially striping the paranthesis from the variable (( ))
		for _, s := range re.FindAllStringSubmatchIndex(value, -1) {

			result = re.ExpandString(result, "$key", value, s)
		}

		resolved := resolver.Get(string(result))

		// if the resolved value is a string return the string
		if v, ok := resolved.(string); ok {
			return v
		}

		// return the orignal value
		return value
	})
}

func evaluateMapStringOfString(template map[string]string, resolver VariableResolver) interface{} {
	transformMap := make(map[string]interface{})
	for k, v := range template {
		transformMap[k] = evaluateString(v, resolver)
	}
	return transformMap
}

func evaluateMapStringOfInterface(template map[string]interface{}, resolver VariableResolver) interface{} {
	for k, v := range template {
		template[k] = evaluate(v, resolver)
	}
	return template
}

func evaluateSliceOfString(template []string, resolver VariableResolver) interface{} {
	// make(type, len, capacity)
	transformSlice := make([]interface{}, 0, len(template))
	for _, v := range template {
		transformSlice = append(transformSlice, evaluateString(v, resolver))
	}
	return transformSlice
}

func evaluateSliceOfInterface(template []interface{}, resolver VariableResolver) interface{} {
	for i, v := range template {
		template[i] = evaluate(v, resolver)
	}
	return template
}
