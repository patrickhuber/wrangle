package filepath

import (
	"path/filepath"
	"strings"
)

// ToSlash replaces all backslashes with forward slashes
func ToSlash(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// Split splits the path into a list
func Split(path string) (dir, file string) {
	slashPath := ToSlash(path)
	dir, file = filepath.Split(slashPath)
	return ToSlash(dir), file
}

// SplitAll splits the path into a list
func SplitAll(path string) []string {
	slashPath := ToSlash(path)
	return strings.Split(slashPath, "/")
}

// Join joins the segments into a single string delimited by forward slash
func Join(segments ...string) string {
	path := strings.Join(segments, "/")
	if len(segments) == 1 {
		if segments[0] == "" {
			path = "/" + path
		}
	}
	return ToSlash(path)
}

// Dir returns the path directory by removing the last element of the path.
// the resultant path will always be forwared slashed
func Dir(path string) string {
	path = filepath.Dir(path)
	return ToSlash(path)
}

// Rel is a passthrough function to filepath.Rel where it calls ToSlash() on the result
func Rel(basepath string, targetpath string) (string, error) {
	r, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return "", err
	}
	return ToSlash(r), nil
}
