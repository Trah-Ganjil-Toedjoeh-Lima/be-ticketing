package util

func ElementDifference(after, before []uint) []uint {
	mb := make(map[uint]struct{}, len(before))
	for _, x := range before {
		mb[x] = struct{}{}
	}
	var diff []uint
	for _, x := range after {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
