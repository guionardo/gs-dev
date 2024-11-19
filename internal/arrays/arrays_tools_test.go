package arrays

import (
	"reflect"
	"strings"
	"testing"
)

func TestArrayHasValue(t *testing.T) {
	tests := []struct {
		name     string
		array    []any
		item     any
		comparer func(j, i any) bool
		want     bool
	}{
		{"Find_ints_with_default_comparer", []any{1, 2, 3, 4, 5}, 3, nil, true},
		{"Find_string_with_case_comparer", []any{"apple", "banana", "cherry"}, "Banana", func(j, i any) bool { return strings.ToLower(j.(string)) == strings.ToLower(i.(string)) }, true},
		{"Find_ints_with_missing_value", []any{1, 2}, 3, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayHasValue(tt.array, tt.item, tt.comparer); got != tt.want {
				t.Errorf("ArrayHasValue() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetArraysDiffs(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 7, 8, 9, 10}
	t.Run("Find difference between two int arrays", func(t *testing.T) {
		diffs := GetArraysDiffs(a, b)
		want := []string{"-4", "-5", "+7", "+8", "+9", "+10"}
		if !reflect.DeepEqual(diffs, want) {
			t.Errorf("GetArraysDiffs() = %v, want %v", diffs, want)
		}
	})
}
