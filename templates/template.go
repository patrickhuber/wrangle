package templates

import (
	"regexp"
)

// Template - contains a go object document that contains placeholders for variables
type Template struct {
	document interface{}
}

// NewTemplate - Creates a new template with the given document parameter
func NewTemplate(document interface{}) *Template {
	return &Template{document: document}
}

// Evaluate - Evaluates the tempalte using the variable resolver for variable lookup
func (template *Template) Evaluate(resolver VariableResolver) interface{} {
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

func evaluateString(template string, resolver VariableResolver) string {
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

func evaluateMapStringOfString(template map[string]string, resolver VariableResolver) map[string]string {
	for k, v := range template {
		template[k] = evaluateString(v, resolver)
	}
	return template
}

func evaluateMapStringOfInterface(template map[string]interface{}, resolver VariableResolver) map[string]interface{} {
	for k, v := range template {
		template[k] = evaluate(v, resolver)
	}
	return template
}

func evaluateSliceOfString(template []string, resolver VariableResolver) []string {
	for i, v := range template {
		template[i] = evaluateString(v, resolver)
	}
	return template
}

func evaluateSliceOfInterface(template []interface{}, resolver VariableResolver) []interface{} {
	for i, v := range template {
		template[i] = evaluate(v, resolver)
	}
	return template
}
