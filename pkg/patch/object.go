package patch

import "reflect"

type ObjectUpdate struct {
	Value map[string]any
}

type ObjectOption func(*ObjectUpdate)

func ObjectSetField(name string, value any) ObjectOption {
	return func(u *ObjectUpdate) {
		u.Value[name] = value
	}
}

func NewObject(options ...ObjectOption) *ObjectUpdate {
	update := &ObjectUpdate{
		Value: make(map[string]any),
	}
	for _, option := range options {
		option(update)
	}
	return update
}

func (u *ObjectUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	stru := reflect.Indirect(val)
	updated := false
	for k, v := range u.Value {
		fieldValue := stru.FieldByName(k)
		fieldType := fieldValue.Type()
		switch fieldType.Kind() {
		case reflect.Ptr:
			updated = updated || u.SetPtr(fieldValue, v)
		default:
			updated = updated || u.Set(fieldValue, v)
		}
	}

	return val, updated
}

func (u *ObjectUpdate) Set(val reflect.Value, value any) bool {
	switch t := value.(type) {
	case Applicable:
		result, ok := t.Apply(val)
		if !ok {
			return false
		}
		val.Set(result)
		return ok
	default:
		if reflect.DeepEqual(val.Interface(), value) {
			return false
		}
		val.Set(reflect.ValueOf(value))
		return true
	}
}

func (u *ObjectUpdate) SetPtr(ptr reflect.Value, value any) bool {
	newValue := reflect.New(ptr.Type().Elem())
	if !u.Set(newValue.Elem(), value) {
		return false
	}
	ptr.Set(newValue)
	return true
}
