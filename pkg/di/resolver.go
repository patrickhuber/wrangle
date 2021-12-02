package di

import "reflect"

type Resolver interface {
	ResolveByType(reflect.Type) interface{}
	Resolve(key string) interface{}
}
