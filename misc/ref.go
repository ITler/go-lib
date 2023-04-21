package misc

// Ref obtains a pointer to the provided value
func Ref[T any](in T) *T {
	return &in
}

// Deref gives back the value of a provided pointer,
// even if in is nil
func Deref[T any](in *T) (result T) {
	if in == nil {
		return
	}
	return *in
}
