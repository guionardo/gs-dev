package arrays

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayHasValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []any
		item     any
		comparer func(j, i any) bool
		want     bool
	}{
		{"Find_ints_with_default_comparer", []any{1, 2, 3, 4, 5}, 3, nil, true},
		{"Find_string_with_case_comparer", []any{"apple", "banana", "cherry"}, "Banana", func(j, i any) bool { return strings.EqualFold(j.(string), i.(string)) }, true},
		{"Find_ints_with_missing_value", []any{1, 2}, 3, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, ArrayHasValue(tt.array, tt.item, tt.comparer), "ArrayHasValue() should return the correct value")
		})
	}
}

func TestGetArraysDiffs(t *testing.T) {
	t.Parallel()

	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 7, 8, 9, 10}

	t.Run("Find difference between two int arrays", func(t *testing.T) {
		t.Parallel()

		diffs := GetArraysDiffs(a, b)

		want := []string{"-4", "-5", "+7", "+8", "+9", "+10"}
		require.Equal(t, want, diffs, "GetArraysDiffs() should return the correct diffs")
	})
}
