package template

import "strings"

func HasVariables(s string) bool {
	return variableRegex.MatchString(s)
}

func ListVariables(s string) []string {
	var vars []string
	matches := variableRegex.FindAllString(s, -1)
	if len(vars) == 0 {
		return vars
	}
	for _, match := range matches {
		v := strings.Trim(match, "()")
		vars = append(vars, v)
	}
	return vars
}
