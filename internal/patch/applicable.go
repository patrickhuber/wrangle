package patch

import "reflect"

type Applicable interface {
	Apply(val reflect.Value) (reflect.Value, bool)
}
