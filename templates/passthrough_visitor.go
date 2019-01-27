package templates

import (
	"fmt"
	"reflect"
)

type passthroughVisitor struct {
}

func NewPassthroughVisitor() Visitor {
	return &passthroughVisitor{}
}

func (visitor *passthroughVisitor) VisitString(value string) (interface{}, error) {
	return value, nil
}

func (visitor *passthroughVisitor) VisitMapStringOfString(value map[string]string) (interface{}, error) {
	transformMap := make(map[string]interface{})
	for k, v := range value {
		newValue, err := visitor.VisitString(v)
		if err != nil {
			return nil, err
		}
		transformMap[k] = newValue
	}
	return transformMap, nil
}

func (visitor *passthroughVisitor) VisitMapStringOfInterface(value map[string]interface{}) (interface{}, error) {
	transformMap := make(map[string]interface{})
	for k, v := range value {
		newValue, err := visitor.VisitInterface(v)
		if err != nil {
			return nil, err
		}
		transformMap[k] = newValue
	}
	return transformMap, nil
}

func (visitor *passthroughVisitor) VisitMapInterfaceOfInterface(value map[interface{}]interface{}) (interface{}, error) {
	transformMap := make(map[interface{}]interface{})
	for k, v := range value {
		newValue, err := visitor.VisitInterface(v)
		if err != nil {
			return nil, err
		}
		transformMap[k] = newValue
	}
	return transformMap, nil
}

func (visitor *passthroughVisitor) VisitInterface(value interface{}) (interface{}, error) {
	switch t := value.(type) {
	case (string):
		return visitor.VisitString(t)
	case (map[string]string):
		return visitor.VisitMapStringOfString(t)
	case (map[string]interface{}):
		return visitor.VisitMapStringOfInterface(t)
	case (map[interface{}]interface{}):
		return visitor.VisitMapInterfaceOfInterface(t)
	case ([]string):
		return visitor.VisitSliceOfString(t)
	case ([]interface{}):
		return visitor.VisitSliceOfInterface(t)
	}

	return nil, fmt.Errorf("unable to evaluate template. Invalid type '%v'", reflect.TypeOf(value))
}

func (visitor *passthroughVisitor) VisitSliceOfString(value []string) (interface{}, error) {
	return value, nil
}

func (visitor *passthroughVisitor) VisitSliceOfInterface(value []interface{}) (interface{}, error) {

	return value, nil
}
