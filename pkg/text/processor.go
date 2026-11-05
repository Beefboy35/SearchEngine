package text

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`[^\w\s]`)

func LowerCase(str string) string {
	return strings.ToLower(str)
}

func Tokenize(str string) []string {
	str = re.ReplaceAllString(str, "")
	return strings.Fields(str)
}
