package templates

import (
	"regexp"
)

type Template struct {
	document interface{}
}

func NewTemplate(document interface{}) *Template {
	return &Template{document: document}
}

func (template *Template) Evaluate(resolver VariableResolver) interface{} {
	return evaluate(template.document, resolver)
}

func evaluate(document interface{}, resolver VariableResolver) interface{} {
	switch t := document.(type) {
	case (string):
		return evaluateString(t, resolver)
	case (map[string]string):
		return evaluateMapStringOfString(t, resolver)
	case (map[string]interface{}):
		return evaluateMapStringOfInterface(t, resolver)
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
