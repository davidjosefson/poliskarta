package filterdescription

import (
	"poliskarta/helperfunctions"
	"strings"
)

var invalidWordsForRoads = []string{
	"Lv", "Länsväg", "länsväg",
	"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65",
}
var validWordsForPlaces = []string{
	"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa",
	"trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj",
	"hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar",
	"dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata",
	"gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata",
	"hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö",
	"träsk", "flygplats", "industriväg", "trappgata", "kärr", "ringvägen",
}

func FilterDescriptionWords(description string) []string {
	locationWords := addWords(description)
	removeInvalidWords(&locationWords)
	return locationWords
}

func addWords(description string) []string {
	prevWordAdded := false

	//Split string on spaces - descWords = array
	descWords := strings.Split(description, " ")

	//Remove spaces from words
	helperfunctions.TrimSpacesFromArray(&descWords)

	//The resulting array of location words after filtering
	locationWords := []string{}

	//Loop through the array of words
	for i := 1; i < len(descWords); i++ {
		currentWord := descWords[i]
		prevWord := descWords[i-1]
		addWord := false

		if currentWord == "" {
			continue
		}

		//Skip iteration if the previous word had a "." in the end
		if strings.HasSuffix(prevWord, ".") {
			prevWordAdded = false
			continue
		}

		//Check if previous word was added and current word is in valid road list
		if prevWordAdded && helperfunctions.TrimmedStringInSlice(currentWord, validWordsForPlaces) {
			addWord = true

			//Add word if it starts with upper case or is a number
		} else if helperfunctions.StartsWithUppercase(currentWord) || helperfunctions.WordIsNumber(currentWord) {
			addWord = true
		}

		//Dont add word if it has been added before
		if helperfunctions.StringInSlice(currentWord, locationWords) {
			addWord = false
		}

		if addWord {
			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")
			locationWords = append(locationWords, currentWord)
			prevWordAdded = true

		} else {
			prevWordAdded = false
		}

	}

	removeInvalidWords(&locationWords)

	return locationWords
}

// func addWords(description string) []string {
// 	prevWordAdded := false

// 	//Split string on spaces - descWords = array
// 	descWords := strings.Split(description, " ")

// 	//Remove spaces from words
// 	helperfunctions.TrimSpacesFromArray(&descWords)

// 	//The resulting array of location words after filtering
// 	locationWords := []string{}

// 	//Loop through the array of words
// 	for i := 1; i < len(descWords); i++ {
// 		currentWord := descWords[i]
// 		prevWord := descWords[i-1]
// 		addWord := false;

// 		if currentWord == "" {
// 			continue
// 		}

// 		//Skip iteration if the previous word had a "." in the end
// 		if strings.HasSuffix(prevWord, ".") {
// 			prevWordAdded = false
// 			continue
// 		}

// 		//Check if previous word was added and current word is in valid road list
// 		if prevWordAdded {
// 			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")

// 			if helperfunctions.StringInSlice(currentWord, validWordsForPlaces) {
// 				locationWords = append(locationWords, currentWord)
// 				prevWordAdded = true
// 				continue
// 			} else {
// 				prevWordAdded = false
// 			}
// 		}

// 		//Check if current word starts with uppercase
// 		if helperfunctions.StartsWithUppercase(currentWord) {
// 			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")
// 			locationWords = append(locationWords, currentWord)
// 			prevWordAdded = true
// 		} else {
// 			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")
// 			//Add word if it is a number
// 			if helperfunctions.WordIsNumber(currentWord) {
// 				locationWords = append(locationWords, currentWord)
// 				prevWordAdded = true
// 			} else {
// 				prevWordAdded = false
// 			}
// 		}

// 		if(addWord)

// 	}

// 	removeInvalidWords(&locationWords)

// 	return locationWords
// }

func removeInvalidWords(locationWords *[]string) {
	sliceCopy := *locationWords

	for index, word := range sliceCopy {
		//Checks if word is invalid ("Länsväg", "E4")
		if helperfunctions.StringInSlice(word, invalidWordsForRoads) {
			//Remove word from slice
			helperfunctions.RemoveIndexFromSlice(index, &sliceCopy)
		} else if
		//Checks case: "E 22" or "E 4" and removes both E and the number
		word == "E" && helperfunctions.CurrentIndexNotLast(index, sliceCopy) && helperfunctions.WordIsNumber(sliceCopy[index+1]) {
			helperfunctions.RemoveIndexFromSlice(index, &sliceCopy)
			helperfunctions.RemoveIndexFromSlice(index+1, &sliceCopy)
		} else if
		//Checks case: "väg 22" and removes the word "väg" and keeps the number
		word == "väg" && helperfunctions.CurrentIndexNotLast(index, sliceCopy) && helperfunctions.WordIsNumber(sliceCopy[index+1]) {
			helperfunctions.RemoveIndexFromSlice(index, &sliceCopy)
		}
	}

	*locationWords = sliceCopy
}
