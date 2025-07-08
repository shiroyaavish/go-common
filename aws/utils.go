package aws

// pointer returns a pointer to any type
func pointer[T any](data T) *T {
	return &data
}
