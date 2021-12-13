package di

import (
	"fmt"
	"reflect"
)

type Container interface {
	RegisterInstance(t reflect.Type, instance interface{})
	RegisterDynamic(t reflect.Type, delegate func(Resolver) (interface{}, error))
	RegisterConstructor(constructor interface{}) error
	Resolver
}

type container struct {
	data map[string]func(Resolver) (interface{}, error)
}

func NewContainer() Container {

	return &container{
		data: map[string]func(Resolver) (interface{}, error){},
	}
}

func (c *container) RegisterConstructor(constructor interface{}) error {
	t := reflect.TypeOf(constructor)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("constructor '%s' must be a method", t.Elem())
	}

	outCount := t.NumOut()
	if outCount == 0 {
		return fmt.Errorf("constructor must have a return value and optional error")
	}
	returnType := t.Out(0)
	if outCount == 2 {
		errorType := t.Out(1)
		if !errorType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("if a constructor has two parameters, the second must implement error")
		}
	} else if outCount != 1 {
		return fmt.Errorf("constructor must have a return value and optional error")
	}

	delegate := func(r Resolver) (interface{}, error) {
		inCount := t.NumIn()
		values := make([]reflect.Value, inCount)
		for i := 0; i < inCount; i++ {
			parameterType := t.In(i)
			parameterFunc, ok := c.data[parameterType.String()]
			if !ok || parameterFunc == nil {
				return nil, fmt.Errorf("error resolving constructor %s missing parameter of type %s", t.String(), parameterType.String())
			}
			value, err := parameterFunc(r)
			if err != nil {
				return nil, err
			}
			values[i] = reflect.ValueOf(value)
		}
		constructorValue := reflect.ValueOf(constructor)
		results := constructorValue.Call(values)
		if len(results) == 0 {
			return nil, fmt.Errorf("no result while executing constructor '%s'", t.String())
		}
		var instance interface{}
		if !results[0].IsNil() {
			instance = results[0].Interface()
		}
		var err error = nil
		if len(results) == 2 {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		}
		return instance, err
	}

	c.data[returnType.String()] = delegate
	return nil
}

func (c *container) RegisterDynamic(t reflect.Type, delegate func(Resolver) (interface{}, error)) {
	c.data[t.String()] = delegate
}

func (c *container) RegisterInstance(t reflect.Type, instance interface{}) {
	c.RegisterDynamic(t, func(r Resolver) (interface{}, error) {
		return instance, nil
	})
}

func (c *container) Resolve(t reflect.Type) (interface{}, error) {
	delegate, ok := c.data[t.String()]
	if !ok {
		return nil, fmt.Errorf("type %s not found", t.String())
	}
	return delegate(c)
}

func (c *container) ResolveAll(t reflect.Type) ([]interface{}, error) {
	return nil, nil
}
