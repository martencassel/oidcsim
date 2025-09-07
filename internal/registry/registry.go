package registry

import "fmt"

type Registry[T any] struct {
	items map[string]T
}

func New[T any]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

func (r *Registry[T]) Register(name string, item T) {
	r.items[name] = item
}

func (r *Registry[T]) Get(name string) (T, error) {
	v, ok := r.items[name]
	if !ok {
		var zero T
		return zero, fmt.Errorf("no item registered with name %q", name)
	}
	return v, nil
}

func (r *Registry[T]) All() map[string]T {
	return r.items
}
