package helperfunctions

import (
	"strings"
	"unicode"
)

func TrimSpacesFromArray(title []string) []string {
	for i := 0; i < len(title); i++ {
		title[i] = strings.TrimSpace(title[i])
	}
	return title
}

func TrimSuffixFromWord(word string, suffix string) string {
	word = strings.TrimSuffix(word, suffix)
	return word
}

func StartsWithUppercase(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

func StringInSliceIgnoreCase(s string, list []string) bool {
	for _, str := range list {
		if strings.ToLower(s) == strings.ToLower(str) {
			return true
		}
	}
	return false
}
