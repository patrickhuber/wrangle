package di

import "reflect"

type Container interface {
	RegisterStaticWithType(t reflect.Type, instance interface{})
	RegisterDynamicWithType(t reflect.Type, delegate func(Resolver) interface{})
	RegisterStatic(key string, instance interface{})
	RegisterDynamic(key string, delegate func(Resolver) interface{})
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

func (c *container) RegisterDynamicWithType(t reflect.Type, delegate func(Resolver) interface{}) {
	c.RegisterDynamic(t.Name(), delegate)
}

func (c *container) RegisterDynamic(key string, delegate func(Resolver) interface{}) {
	c.data[key] = delegate
}

func (c *container) RegisterStaticWithType(t reflect.Type, instance interface{}) {
	c.RegisterStatic(t.Name(), instance)
}

func (c *container) RegisterStatic(key string, instance interface{}) {
	c.RegisterDynamic(key, func(r Resolver) interface{} {
		return instance
	})
}

func (c *container) ResolveByType(t reflect.Type) interface{} {
	delegate, ok := c.data[t.Name()]
	if !ok {
		return nil
	}
	return delegate(c)
}

func (c *container) Resolve(key string) interface{} {
	delegate, ok := c.data[key]
	if !ok {
		return nil
	}
	return delegate(c)
}
