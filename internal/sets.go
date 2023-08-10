package internal

func Set(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, item := range list {
		set[item] = struct{}{}
	}
	return set
}

func List(set map[string]struct{}) []string {
	list := make([]string, len(set))
	index := 0
	for key := range set {
		list[index] = key
		index++
	}
	return list
}
