package patch

import "reflect"

type SliceUpdate struct {
	Patches []*SlicePatch
}

func NewSlice(options ...SliceOption) *SliceUpdate {
	slice := &SliceUpdate{}
	for _, option := range options {
		option(slice)
	}
	return slice
}

type SliceOperation int

const (
	AddOperation     SliceOperation = 0
	RemoveOperation  SliceOperation = 1
	ReplaceOperation SliceOperation = 2
)

type SliceOption = func(*SliceUpdate)

func SliceAppend(value interface{}) SliceOption {
	patch := &SlicePatch{
		Operation: AddOperation,
		Value:     value,
	}
	return func(u *SliceUpdate) {
		u.Patches = append(u.Patches, patch)
	}
}

func SliceRemoveAt(index int) SliceOption {
	patch := &SlicePatch{
		Operation: RemoveOperation,
		Index:     index,
	}
	return func(u *SliceUpdate) {
		u.Patches = append(u.Patches, patch)
	}
}

func SliceRemove(condition func(reflect.Value) bool) SliceOption {
	patch := &SlicePatch{
		Operation: RemoveOperation,
		Condition: condition,
		Index:     -1,
	}
	return func(u *SliceUpdate) {
		u.Patches = append(u.Patches, patch)
	}
}

func SliceModifyAt(index int, value interface{}) SliceOption {
	patch := &SlicePatch{
		Operation: ReplaceOperation,
		Index:     index,
		Value:     value,
	}
	return func(u *SliceUpdate) {
		u.Patches = append(u.Patches, patch)
	}
}

func SliceModify(condition func(reflect.Value) bool, value interface{}) SliceOption {
	patch := &SlicePatch{
		Operation: ReplaceOperation,
		Condition: condition,
		Value:     value,
		Index:     -1,
	}
	return func(u *SliceUpdate) {
		u.Patches = append(u.Patches, patch)
	}
}

type SlicePatch struct {
	Index     int
	Condition func(reflect.Value) bool
	Operation SliceOperation
	Value     interface{}
}

func (u *SliceUpdate) Apply(val reflect.Value) (reflect.Value, bool) {
	var added, removed, modified bool
	var slice reflect.Value
	slice, added = u.Append(val, u.Patches)
	slice, modified = u.Modify(slice, u.Patches)
	slice, removed = u.Filter(slice, u.Patches)
	return slice, added || removed || modified
}

func (u *SliceUpdate) Append(slice reflect.Value, patches []*SlicePatch) (reflect.Value, bool) {
	changed := false

	for _, a := range patches {
		if a.Operation != AddOperation {
			continue
		}
		v := reflect.ValueOf(a.Value)
		slice = reflect.Append(slice, v)
		changed = true
	}
	return slice, changed
}

func (u *SliceUpdate) Filter(slice reflect.Value, patches []*SlicePatch) (reflect.Value, bool) {
	changed := false
	result := reflect.MakeSlice(slice.Type(), 0, 0)
	for i := 0; i < slice.Len(); i++ {
		skip := false

		v := slice.Index(i)
		for _, p := range patches {
			if p.Operation != RemoveOperation {
				continue
			}
			if p.Condition != nil {
				if p.Condition(v) {
					skip = true
				}
			} else if i == p.Index {
				skip = true
			}
		}

		if skip {
			changed = true
			continue
		}
		result = reflect.Append(result, v)
	}
	return result, changed
}

func (u *SliceUpdate) Modify(slice reflect.Value, patches []*SlicePatch) (reflect.Value, bool) {
	changed := false
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		for _, m := range patches {
			if m.Operation != ReplaceOperation {
				continue
			}
			switch t := m.Value.(type) {
			case Applicable:
				applied, modified := t.Apply(v)
				if !modified {
					continue
				}
				v = applied
				changed = true
			default:
				if reflect.DeepEqual(v.Interface(), t) {
					continue
				}
				changed = true
				v = reflect.ValueOf(t)
			}
		}
		slice.Index(i).Set(v)
	}
	return slice, changed
}
