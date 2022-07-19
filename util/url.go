package util

import "net/url"

func UrlMustParse(rawurl string) *url.URL {
	u, _ := url.Parse(rawurl)
	return u
}
