package util

//Contains determines if slice contains value
func Contains(s []int64, i int64) bool {
	for _, v := range s {
		if v == i {
			return true
		}
	}
	return false
}

//Remove an element from a slice - opposite of append
func Remove(s []int64, i int64) []int64 {
	t := []int64{}
	for _, v := range s {
		if v != i {
			t = append(t, v)
		}
	}
	return t
}
