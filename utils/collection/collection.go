package collection

import (
	"reflect"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func NewSetFromSlice[T comparable](slice []T) Set[T] {
	s := make(Set[T])
	for _, v := range slice {
		s.Add(v)
	}
	return s
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for v := range s {
		result = append(result, v)
	}
	return result
}

func (s Set[T]) ForEach(fn func(T)) {
	for v := range s {
		fn(v)
	}
}

func (s Set[T]) Filter(fn func(T) bool) Set[T] {
	result := NewSet[T]()
	for v := range s {
		if fn(v) {
			result.Add(v)
		}
	}
	return result
}

func (s Set[T]) Map(fn func(T) T) Set[T] {
	result := NewSet[T]()
	for v := range s {
		result.Add(fn(v))
	}
	return result
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	result := NewSet[T]()
	for v := range s {
		result.Add(v)
	}
	for v := range other {
		result.Add(v)
	}
	return result
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	result := NewSet[T]()
	for v := range s {
		if other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

func (s Set[T]) Diff(other Set[T]) Set[T] {
	result := NewSet[T]()
	for v := range s {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

func (s Set[T]) IsSubsetOf(other Set[T]) bool {
	for v := range s {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for v := range s {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func SliceUnique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func SliceFilter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func SliceMap[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func SliceGroupBy[T any, K comparable](slice []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range slice {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

func SlicePartition[T any](slice []T, fn func(T) bool) ([]T, []T) {
	yes, no := make([]T, 0), make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return yes, no
}

func SliceIntersect[T comparable](a, b []T) []T {
	setB := NewSetFromSlice(b)
	result := make([]T, 0)
	for _, v := range a {
		if setB.Contains(v) {
			result = append(result, v)
		}
	}
	return result
}

func SliceDiff[T comparable](a, b []T) []T {
	setB := NewSetFromSlice(b)
	result := make([]T, 0)
	for _, v := range a {
		if !setB.Contains(v) {
			result = append(result, v)
		}
	}
	return result
}

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func FilterMap[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

func MapValues[K comparable, V any, R any](m map[K]V, fn func(V) R) map[K]R {
	result := make(map[K]R)
	for k, v := range m {
		result[k] = fn(v)
	}
	return result
}

func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func Clone[T any](src T) T {
	return reflect.ValueOf(src).Interface().(T)
}
