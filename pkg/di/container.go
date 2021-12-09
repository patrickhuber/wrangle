package di

import (
	"fmt"
	"reflect"
)

type Container interface {
	RegisterInstance(t reflect.Type, instance interface{})
	RegisterDynamic(t reflect.Type, delegate func(Resolver) interface{})
	RegisterConstructor(constructor interface{}) error
	Resolver
}

type container struct {
	data map[string]func(Resolver) interface{}
}

func NewContainer() Container {
	return &container{
		data: map[string]func(Resolver) interface{}{},
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
	parameterFunctions := []func(Resolver) interface{}{}
	for i := 0; i < inCount; i++ {
		parameterType := t.In(i)
		parameterFunc := c.data[parameterType.String()]
		parameterFunctions = append(parameterFunctions, parameterFunc)
	}

	delegate := func(r Resolver) interface{} {
		values := []reflect.Value{}
		for _, f := range parameterFunctions {
			value := f(r)
			rv := reflect.ValueOf(value)
			values = append(values, rv)
		}
		result := reflect.ValueOf(constructor).Call(values)
		for _, r := range result {
			return r.Interface()
		}
		return nil
	}

	o := t.Out(0)
	c.data[o.String()] = delegate
	return nil
}

func (c *container) RegisterDynamic(t reflect.Type, delegate func(Resolver) interface{}) {
	c.data[t.String()] = delegate
}

func (c *container) RegisterInstance(t reflect.Type, instance interface{}) {
	c.RegisterDynamic(t, func(r Resolver) interface{} {
		return instance
	})
}

func (c *container) Resolve(t reflect.Type) interface{} {
	delegate, ok := c.data[t.String()]
	if !ok {
		return nil
	}
	return delegate(c)
}
