package internal

func List(set map[string]struct{}) []string {
	list := make([]string, len(set))
	index := 0
	for key := range set {
		list[index] = key
		index++
	}
	return list
}
