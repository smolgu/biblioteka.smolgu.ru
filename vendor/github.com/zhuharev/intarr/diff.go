package intarr

// Diff return removed and new items
func Diff(oldIds, currentIds []int) (joined, leaved []int) {
	var (
		m       = map[int]struct{}{}
		mLeaved = map[int]struct{}{}
	)
	for _, v := range oldIds {
		m[v] = struct{}{}
	}
	for _, v := range currentIds {
		mLeaved[v] = struct{}{}
	}

	for v := range mLeaved {
		if _, ok := m[v]; !ok {
			joined = append(joined, v)
		}
	}

	for v := range m {
		if _, ok := mLeaved[v]; !ok {
			leaved = append(leaved, v)
		}
	}

	return
}
