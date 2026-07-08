package vynvpn

import "net/url"

// makeQuery builds a url.Values from key-value pairs.
func makeQuery(pairs ...string) url.Values {
	q := url.Values{}
	for i := 0; i < len(pairs)-1; i += 2 {
		if pairs[i+1] != "" {
			q.Set(pairs[i], pairs[i+1])
		}
	}
	return q
}
