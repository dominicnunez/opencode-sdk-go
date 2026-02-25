package opencode

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// PtrString returns a pointer to the given string value.
// Useful for optional fields in API request structs.
func PtrString(value string) *string { return &value }

// PtrInt returns a pointer to the given int64 value.
// Useful for optional fields in API request structs.
func PtrInt(value int64) *int64 { return &value }

// PtrFloat returns a pointer to the given float64 value.
// Useful for optional fields in API request structs.
func PtrFloat(value float64) *float64 { return &value }

// PtrBool returns a pointer to the given bool value.
// Useful for optional fields in API request structs.
func PtrBool(value bool) *bool { return &value }
