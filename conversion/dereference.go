package conversion

// DereferenceString safely returns a string value from a string pointer
func DereferenceString(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}
