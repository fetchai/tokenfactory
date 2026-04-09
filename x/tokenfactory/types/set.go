package types

// Set is a generic set implementation using map
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet creates a new set with optional initial values
func NewSet[T comparable](values ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{})}
	s.Add(values...)
	return s
}

// SetFromSlice creates a set from a slice
func SetFromSlice[T comparable](values []T) *Set[T] {
	return NewSet(values...)
}

// Add inserts one or more elements
func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		s.m[v] = struct{}{}
	}
}

// Remove deletes one or more elements
func (s *Set[T]) Remove(values ...T) {
	for _, v := range values {
		delete(s.m, v)
	}
}

// Contains checks if element exists
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len returns number of elements
func (s *Set[T]) Len() int {
	return len(s.m)
}

// IsEmpty checks if set is empty
func (s *Set[T]) IsEmpty() bool {
	return len(s.m) == 0
}

// Clear removes all elements
func (s *Set[T]) Clear() {
	s.m = make(map[T]struct{})
}

// ToSlice returns elements as slice
// Warning: NON-deterministic!
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.m))
	for v := range s.m {
		result = append(result, v)
	}
	return result
}

// Copy creates a shallow copy
func (s *Set[T]) Copy() *Set[T] {
	newSet := NewSet[T]()
	for v := range s.m {
		newSet.m[v] = struct{}{}
	}
	return newSet
}

// Equals checks if two sets are equal
func (s *Set[T]) Equals(other *Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for v := range s.m {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

// IsSubset checks if s ⊆ other
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	for v := range s.m {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

// IsSuperset checks if s ⊇ other
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

// Union returns a new set with elements from both sets
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := s.Copy()
	for v := range other.m {
		result.m[v] = struct{}{}
	}
	return result
}

// Intersection returns common elements
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.m {
		if other.Contains(v) {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// Difference returns elements in s but not in other
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.m {
		if !other.Contains(v) {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// SymmetricDifference returns elements in either set but not both
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.m {
		if !other.Contains(v) {
			result.m[v] = struct{}{}
		}
	}
	for v := range other.m {
		if !s.Contains(v) {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// ForEach iterates over elements
func (s *Set[T]) ForEach(fn func(T)) {
	for v := range s.m {
		fn(v)
	}
}

// Filter returns a new set with elements matching predicate
func (s *Set[T]) Filter(fn func(T) bool) *Set[T] {
	result := NewSet[T]()
	for v := range s.m {
		if fn(v) {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// Any returns true if any element matches predicate
func (s *Set[T]) Any(fn func(T) bool) bool {
	for v := range s.m {
		if fn(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements match predicate
func (s *Set[T]) All(fn func(T) bool) bool {
	for v := range s.m {
		if !fn(v) {
			return false
		}
	}
	return true
}
