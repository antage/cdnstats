package main

import (
	"net/url"
	"strings"
	"fmt"
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

func humanizeSize(b uint64) string {
	if b > 1024*1024*1024*1024 {
		return fmt.Sprintf("%d TiB", b/(1024*1024*1024*1024))
	}
	if b > 1024*1024*1024 {
		return fmt.Sprintf("%d GiB", b/(1024*1024*1024))
	}
	if b > 1024*1024 {
		return fmt.Sprintf("%d MiB", b/(1024*1024))
	}
	if b > 1024 {
		return fmt.Sprintf("%d KiB", b/1024)
	}
	return fmt.Sprintf("%d B", b)
}
