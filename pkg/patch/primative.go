package patch

import (
	"reflect"
	"strings"
)

type StringUpdate struct {
	Value    string
	HasValue bool
}

func NewString(value string) *StringUpdate {
	return &StringUpdate{
		Value:    value,
		HasValue: true,
	}
}

func NewEmptyString() *StringUpdate {
	return &StringUpdate{
		HasValue: false,
	}
}

func (u *StringUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	current := val.String()
	if u.HasValue && !strings.EqualFold(current, u.Value) {
		return reflect.ValueOf(u.Value), true
	}
	return val, false
}

type IntUpdate struct {
	Value    int
	HasValue bool
}

func NewInt(value int) *IntUpdate {
	return &IntUpdate{
		Value:    value,
		HasValue: true,
	}
}

func NewEmptyInt() *IntUpdate {
	return &IntUpdate{
		HasValue: false,
	}
}

func (u *IntUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	current := val.Int()
	if u.HasValue && current != int64(u.Value) {
		return reflect.ValueOf(u.Value), true
	}
	return val, false
}

type BoolUpdate struct {
	Value    bool
	HasValue bool
}

func NewBool(value bool) *BoolUpdate {
	return &BoolUpdate{
		Value:    value,
		HasValue: true,
	}
}

func NewEmptyBool() *BoolUpdate {
	return &BoolUpdate{
		HasValue: false,
	}
}

func (u *BoolUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	current := val.Bool()
	if u.HasValue && current != u.Value {
		return reflect.ValueOf(u.Value), true
	}
	return val, false
}
