package patch

import "reflect"

type SliceUpdate struct {
	Patches []*Patch
}

func NewSlice(patches ...*Patch) *SliceUpdate {
	return &SliceUpdate{
		Patches: patches,
	}
}

type PatchOperation int

const (
	PatchAdd     PatchOperation = 0
	PatchRemove  PatchOperation = 1
	PatchReplace PatchOperation = 2
)

type Patch struct {
	Index     int
	Operation PatchOperation
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

func (u *SliceUpdate) Append(slice reflect.Value, patches []*Patch) (reflect.Value, bool) {
	changed := false

	for _, a := range patches {
		if a.Operation != PatchAdd {
			continue
		}
		v := reflect.ValueOf(a.Value)
		slice = reflect.Append(slice, v)
		changed = true
	}
	return slice, changed
}

func (u *SliceUpdate) Filter(slice reflect.Value, patches []*Patch) (reflect.Value, bool) {
	changed := false
	result := reflect.MakeSlice(slice.Type(), 0, 0)
	for i := 0; i < slice.Len(); i++ {
		skip := false
		for _, p := range patches {
			if p.Operation != PatchRemove {
				continue
			}
			if i == p.Index {
				skip = true
			}
		}
		if skip {
			changed = true
			continue
		}
		v := slice.Index(i)
		result = reflect.Append(result, v)
	}
	return result, changed
}

func (u *SliceUpdate) Modify(slice reflect.Value, patches []*Patch) (reflect.Value, bool) {
	changed := false
	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		for _, m := range patches {
			if m.Operation != PatchReplace {
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
