package internal

import (
	"fmt"
	"testing"
)

func TestSet_Create(t *testing.T) {

	t.Run("Test", func(t *testing.T) {
		set := NewSet()
		if set.Len() != 0 {
			t.Errorf("Set.Len() = %v, want %v", set.Len(), 0)
		}
		set.Add("Sample1")
		set.Add("Sample2")
		set.Add("Sample1")
		if set.Len() != 2 {
			t.Errorf("Set.Len() = %v, want %v", set.Len(), 2)
		}
		set.Remove("Sample1")
		if set.Len() != 1 {
			t.Errorf("Set.Len() = %v, want %v", set.Len(), 1)
		}

		items := set.ToList()
		expected := [...]string{
			"Sample2",
		}
		if fmt.Sprintf("%v", items) != fmt.Sprintf("%v", expected) {
			t.Errorf("Set.ToList() = %v, want %v", items, expected)
		}

	})

}
