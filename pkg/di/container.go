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
	if outCount != 1 {
		return fmt.Errorf("constructor must return one argument")
	}

	inCount := t.NumIn()
	parameterFunctions := map[string]func(Resolver) (interface{}, error){}
	for i := 0; i < inCount; i++ {
		parameterType := t.In(i)
		parameterFunc := c.data[parameterType.String()]
		parameterFunctions[parameterType.String()] = parameterFunc
	}

	delegate := func(r Resolver) (interface{}, error) {
		values := []reflect.Value{}
		for k, f := range parameterFunctions {
			if f == nil {
				return nil, fmt.Errorf("error resolving constructor %s missing parameter of type %s", t.String(), k)
			}
			value, err := f(r)
			if err != nil {
				return nil, err
			}
			rv := reflect.ValueOf(value)
			values = append(values, rv)
		}
		constructorValue := reflect.ValueOf(constructor)
		result := constructorValue.Call(values)
		for _, r := range result {
			return r.Interface(), nil
		}
		return nil, nil
	}

	o := t.Out(0)
	c.data[o.String()] = delegate
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