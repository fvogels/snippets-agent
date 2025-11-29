package util

type Set[T comparable] struct {
	items map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		items: make(map[T]bool),
	}
}

func NewSetFromSlice[T comparable](items []T) Set[T] {
	set := NewSet[T]()

	for _, item := range items {
		set.Add(item)
	}

	return set
}

func (set *Set[T]) Add(item T) {
	set.items[item] = true
}

func (set *Set[T]) Contains(item T) bool {
	_, ok := set.items[item]

	return ok
}

func (set *Set[T]) ToSlice() []T {
	result := make([]T, len(set.items))
	index := 0

	for item := range set.items {
		result[index] = item
		index++
	}

	return result
}

func (set *Set[T]) Intersect(keep Set[T]) {
	toBeRemoved := []T{}

	for _, item := range set.ToSlice() {
		if !keep.Contains(item) {
			toBeRemoved = append(toBeRemoved, item)
		}
	}

	for _, item := range toBeRemoved {
		delete(set.items, item)
	}
}

func (set *Set[T]) Union(other Set[T]) {
	for item := range other.items {
		set.Add(item)
	}
}

func (set *Set[T]) IntersectMany(keep ...Set[T]) {
	for _, s := range keep {
		set.Intersect(s)
	}
}

func (set Set[T]) IsSubsetOf(other Set[T]) bool {
	for item := range set.items {
		if !other.Contains(item) {
			return false
		}
	}

	return true
}

func (set Set[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(set)
}
