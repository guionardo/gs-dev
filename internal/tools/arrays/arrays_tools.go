package arrays

import (
	"fmt"
	"slices"
)

// GetArraysDiffs returns the differences between two slices a and b.
// It returns a slice of strings where each string represents an element
// that is present in one slice but not the other. Elements from slice a
// that are not in slice b are prefixed with "-", and elements from slice b
// that are not in slice a are prefixed with "+".
//
// T is a type parameter that must be comparable.
//
// Parameters:
//
//	a - The first slice to compare.
//	b - The second slice to compare.
//
// Returns:
//
//	A slice of strings representing the differences between the two slices.
func GetArraysDiffs[T comparable](a []T, b []T) []string {
	results := make([]string, 0, len(a)+len(b))
	for i := range a {
		if !slices.Contains(b, a[i]) {
			results = append(results, fmt.Sprintf("-%v", a[i]))
		}
	}

	for i := range b {
		if !slices.Contains(a, b[i]) {
			results = append(results, fmt.Sprintf("+%v", b[i]))
		}
	}

	return results
}
