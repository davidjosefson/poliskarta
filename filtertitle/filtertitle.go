package filtertitle

import (
	"errors"
	"poliskarta/helperfunctions"
	"strings"
)

func FilterTitleWords(title string) ([]string, error) {
	var locationWords []string
	var err error

	if !hasSammanfattning(title) {
		locationWords = strings.Split(title, ",")
		removeTimeStamp(&locationWords)
		trimSpecialChars(&locationWords)
		removeNonLocationWords(&locationWords)
	} else {
		err = errors.New("Titeln är av typen 'Sammanfattning' och innehåller ingen platsinformation")
	}

	return locationWords, err
}

func removeTimeStamp(title *[]string) {
	titleCopy := *title
	helperfunctions.RemoveIndexFromSlice(0, &titleCopy)
	*title = titleCopy
}

func trimSpecialChars(title *[]string) {
	titleCopy := *title

	helperfunctions.TrimPrefixesFromStringSlice(&titleCopy, " ", ",", ".")
	helperfunctions.TrimSuffixesFromStringSlice(&titleCopy, " ", ",", ".")

	*title = titleCopy
}

func removeNonLocationWords(title *[]string) {
	titleCopy := *title
	titleLocation := make([]string, 0)

	//Add all words after the first uppercase
	for i := 1; i < len(titleCopy); i++ {
		if helperfunctions.StartsWithUppercase(titleCopy[i]) {
			for j := i; j < len(titleCopy); j++ {
				titleLocation = append(titleLocation, titleCopy[j])
			}
			break
		}
	}

	*title = titleLocation
}

func hasSammanfattning(title string) bool {
	return strings.Contains(strings.ToLower(title), "sammanfattning")
}
