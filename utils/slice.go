package utils

func IsInSlice[K comparable](item K, items []K) bool {
	for _, inItem := range items {
		if inItem == item {
			return true
		}
	}
	return false
}
