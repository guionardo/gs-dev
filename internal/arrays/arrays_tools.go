package arrays

import "fmt"

// ArrayHasValue checks if a given value exists in an array using one or more comparison functions.
// If no comparison function is provided, it defaults to using the equality operator.
//
// Parameters:
//
//	a        - The array to search within.
//	value    - The value to search for.
//	comparers - Optional comparison functions to use for checking equality.
//
// Returns:
//
//	bool - True if the value is found in the array, otherwise false.
func ArrayHasValue[T comparable](a []T, value T, comparers ...func(j, i T) bool) bool {
	if len(comparers) == 0 || comparers[0] == nil {
		comparers = []func(j, i T) bool{func(j, i T) bool { return j == i }}
	}

	for index := range a {
		for comp := range comparers {
			if comparers[comp](a[index], value) {
				return true
			}
		}
	}

	return false
}

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
		if !ArrayHasValue(b, a[i]) {
			results = append(results, fmt.Sprintf("-%v", a[i]))
		}
	}

	for i := range b {
		if !ArrayHasValue(a, b[i]) {
			results = append(results, fmt.Sprintf("+%v", b[i]))
		}
	}

	return results
}
