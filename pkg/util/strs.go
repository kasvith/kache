package util

// ToString Convert an interface to string
func ToString(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}

	return ""
}
