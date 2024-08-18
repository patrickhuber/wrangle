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

type EvaluationResult struct {
	Value      any
	Unresolved []string
}

func (e *Evaluator) Evaluate(data any) (*EvaluationResult, error) {
	v := reflect.ValueOf(data)
	return e.walk(v)
}

func (e *Evaluator) walk(v reflect.Value) (*EvaluationResult, error) {

	// dereference interfaces
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return e.EvaluateString(v.String())
	case reflect.Int, reflect.Int64, reflect.Float64:
		return &EvaluationResult{Value: v.Interface()}, nil
	case reflect.Slice:
		return e.EvaluateSlice(v)
	case reflect.Map:
		return e.EvaluateMap(v)
	}

	return nil, fmt.Errorf("unsupported type '%s'", v.Kind())
}

func (e *Evaluator) EvaluateString(s string) (*EvaluationResult, error) {
	vars := variableRegex.FindAllString(s, -1)

	if len(vars) == 0 {
		return &EvaluationResult{Value: s}, nil
	}

	if e.variableMap == nil {
		e.variableMap = map[string]any{}
	}

	var unresolved []string

	// add unique
	for _, variable := range vars {

		// trim (( )), cache the original for replacement
		original := variable
		variable = strings.Trim(variable, "()")

		// fetch the variable value
		value, found, err := e.getValue(variable)
		if err != nil {
			return nil, err
		}

		if !found {
			unresolved = append(unresolved, variable)
			continue
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

	// return error of unresolved variables
	if len(unresolved) > 0 {
		return &EvaluationResult{Value: nil, Unresolved: unresolved}, nil
	}

	return &EvaluationResult{Value: s}, nil
}

func (e *Evaluator) getValue(s string) (any, bool, error) {

	// in the cache?
	v, ok := e.variableMap[s]
	if ok {
		return v, true, nil
	}

	var value any
	found := false
	// query all the providers in order to cascade
	for _, p := range e.providers {
		v, ok, err := p.Get(s)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			continue
		}
		value = v
		found = true
	}

	if !found {
		return nil, false, nil
	}

	// set cache
	e.variableMap[s] = value

	return value, true, nil
}

func (e *Evaluator) EvaluateSlice(v reflect.Value) (*EvaluationResult, error) {
	var slice []any
	var unresolved []string
	for i := 0; i < v.Len(); i++ {
		value, err := e.walk(v.Index(i))
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, fmt.Errorf("unexpected nil value from slice evaluation")
		}
		if len(value.Unresolved) > 0 {
			unresolved = append(unresolved, value.Unresolved...)
		}
		slice = append(slice, value.Value)
	}
	return &EvaluationResult{Value: slice, Unresolved: unresolved}, nil
}

func (e *Evaluator) EvaluateMap(v reflect.Value) (*EvaluationResult, error) {
	clone := reflect.MakeMap(v.Type())
	keys := v.MapKeys()

	var unresolved []string
	for _, key := range keys {
		value := v.MapIndex(key)
		newKey, err := e.walk(key)
		if err != nil {
			return nil, err
		}
		if len(newKey.Unresolved) > 0 {
			unresolved = append(unresolved, newKey.Unresolved...)
		}
		newValue, err := e.walk(value)
		if err != nil {
			return nil, err
		}
		if len(newValue.Unresolved) > 0 {
			unresolved = append(unresolved, newValue.Unresolved...)
		}
		clone.SetMapIndex(reflect.ValueOf(newKey.Value), reflect.ValueOf(newValue.Value))
	}
	return &EvaluationResult{Value: clone.Interface(), Unresolved: unresolved}, nil
}
