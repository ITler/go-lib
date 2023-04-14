package ref

// String obtains a pointer to the provided string.
func String(in string) *string {
	return &in
}
