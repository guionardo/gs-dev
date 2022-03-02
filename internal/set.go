package internal

type (
	void struct{}
	Set  struct {
		items map[string]void
	}	
)

var member void

func NewSet() Set {
	return Set{
		items: make(map[string]void),
	}
}

func (set *Set) Len() int {
	return len(set.items)
}

func (set *Set) Add(value string) *Set {
	set.items[value] = member
	return set
}

func (set *Set) Remove(value string) *Set {
	delete(set.items, value)
	return set
}

func (set *Set) ToList() []string {
	list := make([]string, len(set.items))
	index := 0
	for key := range set.items {
		list[index] = key
		index++
	}
	return list
}

