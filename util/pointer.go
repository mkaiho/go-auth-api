package util

func ToPointer[T any](v T) *T {
	return &v
}

func Deref[T any](p *T, defaultValue T) T {
	if p != nil {
		return *p
	}
	return defaultValue
}
