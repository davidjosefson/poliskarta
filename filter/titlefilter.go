package filter

import (
	"errors"
	"poliskarta/helperfunctions"
	"strings"
)

func FilterTitleWords(title string) ([]string, error) {
	var locationWords []string
	var err error

	if HasLocationInTitle(title) {
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

func HasLocationInTitle(title string) bool {
	if strings.Contains(strings.ToLower(title), "sammanfattning") || strings.Contains(strings.ToLower(title), "övrigt") || strings.Contains(strings.ToLower(title), "obemannat") {
		return false
	} else {
		return true
	}
}

func GetTime(title string) string {
	titleWords := strings.SplitN(title, ",", 2)

	//First word is always time
	return titleWords[0]
}

func GetEventType(title string) string {
	titleWords := strings.Split(title, ",")

	eventType := ""
	//Skip first word, which always is time, and last word which is location-info
	for i := 1; i < len(titleWords)-1; i++ {
		trimmedWord := strings.TrimSuffix(titleWords[i], ",")
		trimmedWord = strings.TrimSpace(trimmedWord)
		eventType += trimmedWord + " "
	}

	eventType = strings.TrimSpace(eventType)

	return eventType
}
