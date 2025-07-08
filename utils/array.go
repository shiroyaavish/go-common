package utils

// IsInArray checks if the given data is present in the given array.
// It iterates through the array and compares each element with the data using the == operator.
// If a matching element is found, it returns true. Otherwise, it returns false.
func IsInArray[T comparable](findThisData T, inThisArray []T) bool {
	if inThisArray == nil {
		return false
	}
	for i := 0; i < len(inThisArray); i++ {
		if findThisData == inThisArray[i] {
			return true
		}
	}
	return false
}

// IsInArrayIndex checks if the given data is present in the given array.
// It iterates through the array and compares each element with the data using the == operator.
// If a matching element is found, it returns the index of that element. Otherwise, it returns -1.
func IsInArrayIndex[T comparable](findThisData T, inThisArray []T) int {
	if inThisArray == nil {
		return -1
	}
	for i := 0; i < len(inThisArray); i++ {
		if findThisData == inThisArray[i] {
			return i
		}
	}
	return -1
}
