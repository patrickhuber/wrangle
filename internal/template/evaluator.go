package template

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var variableRegex = regexp.MustCompile(`\(\((!?[-/\.\w\pL]+)\)\)`)

type Evaluator struct {
	providers   []VariableProvider
	variableMap map[string]any
}

func (e *Evaluator) Evaluate(data any) (any, error) {
	v := reflect.ValueOf(data)
	return e.walk(v)
}

func (e *Evaluator) walk(v reflect.Value) (any, error) {

	// dereference interfaces
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return e.EvaluateString(v.String())
	case reflect.Int, reflect.Int64, reflect.Float64:
		return v.Interface(), nil
	case reflect.Slice:
		return e.EvaluateSlice(v)
	case reflect.Map:
		return e.EvaluateMap(v)
	}

	return nil, fmt.Errorf("unsupported type '%s'", v.Kind())
}

func (e *Evaluator) EvaluateString(s string) (any, error) {
	vars := variableRegex.FindAllString(s, -1)

	if len(vars) == 0 {
		return s, nil
	}

	if e.variableMap == nil {
		e.variableMap = map[string]any{}
	}

	// add unique
	for _, variable := range vars {

		// trim (( )), cache the original for replacement
		original := variable
		variable = strings.Trim(variable, "()")

		// fetch the variable value
		value, err := e.getValue(variable)
		if err != nil {
			return nil, err
		}

		// check the type of the result to ensure it is a primitive type
		switch value.(type) {
		case string, int, int16, int32, int64, uint, uint16, uint32, uint64:
			valueString := fmt.Sprintf("%v", value)
			s = strings.ReplaceAll(s, original, valueString)
		default:
			return nil, fmt.Errorf("invalid type '%T' for value '%v' variable '%s'", value, value, variable)
		}
	}

	return s, nil
}

func (e *Evaluator) getValue(s string) (any, error) {

	// in the cache?
	v, ok := e.variableMap[s]
	if ok {
		return v, nil
	}

	var value any
	found := false
	// query all the providers in order to cascade
	for _, p := range e.providers {
		v, ok, err := p.Get(s)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		value = v
		found = true
	}

	if !found {
		return nil, fmt.Errorf("unable to find variable '%s' in any of the configured stores", s)
	}
	// set cache
	e.variableMap[s] = value

	return value, nil
}

func (e *Evaluator) EvaluateSlice(v reflect.Value) (any, error) {
	var slice []any
	for i := 0; i < v.Len(); i++ {
		value, err := e.walk(v.Index(i))
		if err != nil {
			return nil, err
		}
		slice = append(slice, value)
	}
	return slice, nil
}

func (e *Evaluator) EvaluateMap(v reflect.Value) (any, error) {
	clone := reflect.MakeMap(v.Type())
	keys := v.MapKeys()
	for _, key := range keys {
		value := v.MapIndex(key)
		newKey, err := e.walk(key)
		if err != nil {
			return nil, err
		}
		newValue, err := e.walk(value)
		if err != nil {
			return nil, err
		}
		clone.SetMapIndex(reflect.ValueOf(newKey), reflect.ValueOf(newValue))
	}
	return clone.Interface(), nil
}
