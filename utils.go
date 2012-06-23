package main

import (
	"net/url"
	"strings"
)

func normalizeReferer(r string) string {
	if len(r) == 0 {
		return r
	} else {
		result := ""
		u, err := url.Parse(r)
		if err == nil {
			result = u.Host
			if strings.HasPrefix(u.Host, "www.") {
				result = u.Host[4:]
			}
		}
		return result
	}
	panic("unreachable")
}
