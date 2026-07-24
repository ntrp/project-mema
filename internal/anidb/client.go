package anidb

// ValidClientName reports whether name is an AniDB registered-client identifier.
func ValidClientName(name string) bool {
	if name == "" {
		return false
	}
	for _, character := range name {
		if (character < 'a' || character > 'z') && (character < '0' || character > '9') {
			return false
		}
	}
	return true
}
