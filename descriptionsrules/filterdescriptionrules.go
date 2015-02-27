package main

import (
	"poliskarta/helperfunctions"
	"strings"
)

//Used to filter out Europe road names, like "E6", "E22"
var europeRoads []string

// var arrayOfTitles []string

//Used to include common words in names of places
//separated by space, like "Jörgens trappa", "Anna Lindhs plats"
var validWordsForPlaces []string

//Used to filter out words for roads followed by numbers,
//like "Lv 598", "väg 112" and the like
var invalidWordsForRoads []string

//Rule 1:
func Rule1(description string) []string {
	fillEuropeRoads()
	fillValidWordsForPlaces()
	fillInvalidWordsForRoads()

	//Split on spaces - descWords = array
	descWords := strings.Split(description, " ")
	trimmedDescWords := helperfunctions.TrimSpacesFromArray(descWords)

	//The resulting array of words after filtering
	placeWords := []string{}

	//Loop through the array of words
	for i := 1; i < len(trimmedDescWords); i++ {

		currentWord := trimmedDescWords[i]

		//Go to the next word in the array if the previous word had a "." in the end
		if strings.HasSuffix(trimmedDescWords[i-1], ".") {
			continue
		}

		//Check if current word is part of the invalid road-words
		if helperfunctions.StringInSlice(currentWord, invalidWordsForRoads) && currentWordNotLastWordInArray(i, trimmedDescWords) {
			nextWordInArray := trimmedDescWords[i+1]
			helperfunctions.TrimSuffixesFromWord(&nextWordInArray, ".", ",")

			//Check if next word is number, if so: add it
			if helperfunctions.WordIsNumber(nextWordInArray) {
				placeWords = append(placeWords, nextWordInArray)
				continue
			}
		}

		//Check if the word starts with uppercase and is NOT a europe road
		if helperfunctions.StartsWithUppercase(currentWord) && !helperfunctions.StringInSliceIgnoreCase(currentWord, europeRoads) {
			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")
			placeWords = append(placeWords, currentWord)
		}
	}

	return placeWords
}
func fillEuropeRoads() {
	europeRoads = []string{"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65", "E", "Lv"}
}

func fillValidWordsForPlaces() {
	validWordsForPlaces = []string{"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa", "trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj", "hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar", "dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata", "gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata", "hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö", "träsk", "flygplats", "industriväg", "trappgata", "kärr"}
}

func fillInvalidWordsForRoads() {
	invalidWordsForRoads = []string{"väg", "Lv", "Länsväg", "länsväg"}
}

func currentWordNotLastWordInArray(index int, strings []string) bool {
	return index < len(strings)-1
}
