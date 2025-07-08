package utils

func Pointer[T any](data T) *T {
	return &data
}
