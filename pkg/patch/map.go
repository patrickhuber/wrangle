package patch

import "reflect"

type MapUpdate struct {
	Set    map[string]interface{}
	Remove []string
}

type MapOption func(*MapUpdate)

func MapRemove(key string) MapOption {
	return func(m *MapUpdate) {
		if m.Remove == nil {
			m.Remove = []string{}
		}
		m.Remove = append(m.Remove, key)
	}
}

func MapSet(key string, value interface{}) MapOption {
	return func(m *MapUpdate) {
		if m.Set == nil {
			m.Set = map[string]interface{}{}
		}
		m.Set[key] = value
	}
}

func NewMap(options ...MapOption) *MapUpdate {
	update := &MapUpdate{}
	for _, option := range options {
		option(update)
	}
	return update
}

func (u *MapUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	var set, removed bool
	var m reflect.Value
	m, set = u.Modify(val)
	m, removed = u.Filter(m)
	return m, set || removed
}

func (u *MapUpdate) Modify(val reflect.Value) (reflect.Value, bool) {
	changed := false
	for k, v := range u.Set {
		key := reflect.ValueOf(k)
		value := reflect.ValueOf(v)
		current := val.MapIndex(key)

		// is missing?
		if reflect.DeepEqual(current, reflect.Value{}) {
			val.SetMapIndex(key, value)
			changed = true
			continue
		}

		switch t := v.(type) {
		case Applicable:
			applied, modified := t.Apply(current)
			if !modified {
				continue
			}
			value = applied
		default:
			if reflect.DeepEqual(current.Interface(), t) {
				continue
			}
		}

		changed = true
		val.SetMapIndex(key, value)
	}
	return val, changed
}

func (u *MapUpdate) Filter(val reflect.Value) (reflect.Value, bool) {
	changed := false
	for _, r := range u.Remove {
		val.SetMapIndex(reflect.ValueOf(r), reflect.Value{})
		changed = true
	}
	return val, changed
}
