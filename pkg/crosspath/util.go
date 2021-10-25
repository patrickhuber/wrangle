// Package crosspath anonicalizes all paths to forward slash and has a similar signature to the filepath package.
// When developing on windows, the default functions will use backslash which can mess up unit tests.
package crosspath

import "strings"

// Join joins the segments into a single string delimited by forward slash
func Join(segments ...string) string {
	path := ""
	for i, s := range segments {
		s = strings.Replace(s, "\\", "/", -1)
		if i > 0 {
			path = path + "/" + strings.Trim(s, "/")
		} else {
			path = strings.TrimRight(s, "/")
		}
	}
	return path
}

// ToSlash replaces all backslashes with forward slashes
func ToSlash(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
