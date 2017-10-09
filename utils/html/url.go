package html

import (
	"regexp"
	"net/url"
)

var (
	urlAnchorExp *regexp.Regexp
)

func init() {
	urlAnchorExp = regexp.MustCompile("#.*$")
}

func RemoveUrlAnchor(url string) string {
	return urlAnchorExp.ReplaceAllString(url, "")
}

func GetDomain(rawUrl string) string {
	u, _ := url.Parse(rawUrl)
	return u.Host
}
