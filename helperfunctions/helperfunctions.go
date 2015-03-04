package helperfunctions

import (
	"strconv"
	"strings"
	"unicode"
)

func TrimSpacesFromArray(title *[]string) {
	copy := *title

	for i := 0; i < len(*title); i++ {
		copy[i] = strings.TrimSpace(copy[i])
	}

	*title = copy
}

func TrimSuffixesFromWord(word *string, suffixes ...string) {
	for _, suffix := range suffixes {
		copy := *word
		*word = strings.TrimSuffix(copy, suffix)
	}
}

func StartsWithUppercase(str string) bool {
	return unicode.IsUpper([]rune(str)[0])
}

func StringInSlice(str string, slice []string) bool {
	for _, strInSlice := range slice {
		if str == strInSlice {
			return true
		}
	}
	return false
}

func StringInSliceIgnoreCase(str string, slice []string) bool {
	for _, strInSlice := range slice {
		if strings.ToLower(str) == strings.ToLower(strInSlice) {
			return true
		}
	}
	return false
}

func WordIsNumber(word string) bool {
	_, err := strconv.Atoi(word)
	return err == nil
}

func RemoveIndexFromSlice(index int, slice *[]string) {
	sliceCopy := *slice
	sliceCopy = append(sliceCopy[:index], sliceCopy[index+1:]...)
	*slice = sliceCopy
}

func CurrentIndexNotLast(index int, slice []string) bool {
	return index < len(slice)-1
}
