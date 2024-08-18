package stores

import (
	"fmt"
)

type registry struct {
	factories []Factory
}

// Registry provides a factory store registration interface
type Registry interface {
	// Register searches for an existing factory matching the name and registers the factory if none exists
	Register(f Factory)
	// Find attempts to find the factory and returns false if one can not be found
	Find(name string) (Factory, bool)
	// Get attempts to find the factory and returns an error if one can not be found
	Get(name string) (Factory, error)
	// Removes the given item by name
	Remove(name string)
}

// NewRegistry returns a new registry
func NewRegistry(factories []Factory) Registry {
	return &registry{
		factories: factories,
	}
}

func (r *registry) Register(f Factory) {
	f, ok := r.Find(f.Name())
	if !ok {
		r.factories = append(r.factories, f)
	}
}

func (r registry) Find(name string) (Factory, bool) {
	f, i := r.FindIndex(name)
	return f, i != -1
}

func (r registry) FindIndex(name string) (Factory, int) {
	for i, f := range r.factories {
		if f.Name() == name {
			return f, i
		}
	}
	return nil, -1
}

func (r registry) Get(name string) (Factory, error) {
	f, ok := r.Find(name)
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("unable to find factory with name '%s'", name)
}

func (r *registry) Remove(name string) {
	_, i := r.FindIndex(name)
	if i == -1 {
		return
	}
	r.RemoveAt(i)
}

func (r *registry) RemoveAt(index int) {
	var zero Factory
	copy(r.factories[index:], r.factories[index+1:]) // Shift a[i+1:] left one index.
	r.factories[len(r.factories)-1] = zero           // Erase last element (write zero value).
	r.factories = r.factories[:len(r.factories)-1]   // Truncate slice.
}
