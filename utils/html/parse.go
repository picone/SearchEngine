package html

import (
	"regexp"
	"strings"
	"log"
)

var titleExp, metaExp *regexp.Regexp

func init() {
	titleExp = regexp.MustCompile("(?i)<title>([\\w\\W]+)</title>")
	metaExp = regexp.MustCompile("(?i)<meta[\\s\\S]name=\"([\\w\\W]+?)\"[\\s\\S]content=\"([\\w\\W]+?)\"")
}

func ParseTitle(page string) string {
	result := titleExp.FindStringSubmatch(page)
	if len(result) == 2 {
		return result[1]
	} else {
		return ""
	}
}

func ParseMeta(page string) map[string]string {
	result := make(map[string]string)
	matches := metaExp.FindAllStringSubmatch(page, -1)
	if matches != nil {
		for _, match := range matches {
			result[strings.ToLower(match[1])] = match[2]
		}
	}
	log.Println(result)
	return result
}
