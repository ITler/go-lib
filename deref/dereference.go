package deref

// String safely returns a string value from a string pointer
func String(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}
