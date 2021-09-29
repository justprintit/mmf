package util

import (
	"net/url"
	"strings"
)

func NextInPathUnescaped(path string, prefixes ...string) (string, string, int) {
	p0, s, n := NextInPath(path, prefixes...)
	if len(s) > 0 {
		if s, err := url.PathUnescape(s); err == nil {
			return p0, s, n
		}
	}
	return "", "", -1
}

func NextInPath(path string, prefixes ...string) (string, string, int) {
	for _, p0 := range prefixes {
		if p1 := strings.TrimPrefix(path, p0); p1 != path {
			// prefix match
			if n := strings.IndexRune(p1, '/'); n >= 0 {
				p1 = p1[:n]
			}

			n0 := len(p0)
			n1 := len(p1)

			return p0, p1, n0 + n1
		}
	}

	return "", "", -1
}
