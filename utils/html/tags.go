package html

import (
	"regexp"
)

var (
	tagExp, styleExp, scriptExp, spaceExp, hrefExp, specialTagExp *regexp.Regexp
)

func init() {
	tagExp = regexp.MustCompile("<[\\S\\s]+?>")
	styleExp = regexp.MustCompile("(?i)<style[\\S\\s]+?</style>")
	scriptExp = regexp.MustCompile("(?)<script[\\S\\s]+?</script>")
	spaceExp = regexp.MustCompile("\\s{2,}")
	hrefExp = regexp.MustCompile("(?i)<a[\\S\\s]href=\"(http.+?)\"")
	specialTagExp = regexp.MustCompile("\\\\[rntfv]")
}

func RemoveHTMLTags(content string) string {
	result := styleExp.ReplaceAllString(content, "")
	result = scriptExp.ReplaceAllString(result, "")
	result = tagExp.ReplaceAllString(result, "")
	result = specialTagExp.ReplaceAllString(result, " ")
	return spaceExp.ReplaceAllString(result, "")
}

func GetHrefLinks(content string) []string {
	matches := hrefExp.FindAllStringSubmatch(content, -1)
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match[1]
	}
	return result
}
