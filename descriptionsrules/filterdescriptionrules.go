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

	prevWordAdded := false

	//Split on spaces - descWords = array
	descWords := strings.Split(description, " ")
	helperfunctions.TrimSpacesFromArray(&descWords)

	//The resulting array of words after filtering
	locationWords := []string{}

	//Loop through the array of words
	for i := 1; i < len(descWords); i++ {

		currentWord := descWords[i]
		prevWord := descWords[i-1]

		//Skip iteration if the previous word had a "." in the end
		if strings.HasSuffix(prevWord, ".") {
			prevWordAdded = false
			continue
		}

		//Check if previous word was added and current word is in valid road list
		if prevWordAdded {
			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")

			if helperfunctions.StringInSlice(currentWord, validWordsForPlaces) {
				locationWords = append(locationWords, currentWord)
				prevWordAdded = true
				continue
			}
		}

		/*
			Här har jag en tanke om att fixa en bakåt-koll när man stöter på en siffra, istället
			för en framåtkoll. Nu kollas ifall VARJE ord finns med i invalid road-listan!
			Då följer man samma struktur som för övriga checkar vi gör. Men det kan också bli kass, har
			inte tänkt igenom det helt.

			Dessutom: kanske kan fixa en punkt-koll så att man bara behöver göra trim på . och , en gång i början,
			alltså först kolla om ordet har punkt eller komma, och isf sätta en bool till true, sen trimma oavsett.
			Nu trimmas det tre gånger, och fler lär det bli.
		*/

		//Check if current word is part of the invalid road-words
		if helperfunctions.StringInSlice(currentWord, invalidWordsForRoads) && currentIndexNotLast(i, descWords) {
			nextWordInArray := descWords[i+1]
			helperfunctions.TrimSuffixesFromWord(&nextWordInArray, ".", ",")

			//Check if next word is number, if so: add it
			if helperfunctions.WordIsNumber(nextWordInArray) {
				locationWords = append(locationWords, nextWordInArray)
				prevWordAdded = true
				i++
				continue
			}
		}

		//Check if current word starts with uppercase and is NOT europe road
		if helperfunctions.StartsWithUppercase(currentWord) {
			helperfunctions.TrimSuffixesFromWord(&currentWord, ".", ",")
			if !helperfunctions.StringInSliceIgnoreCase(currentWord, europeRoads) {
				addWordTolocationWords()
				locationWords = append(locationWords, currentWord)
				prevWordAdded = true
			}
		} else {
			prevWordAdded = false
		}
	}

	return locationWords
}
func fillEuropeRoads() {
	europeRoads = []string{"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65", "E", "Lv"}
}

func fillValidWordsForPlaces() {
	validWordsForPlaces = []string{"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa", "trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj", "hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar", "dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata", "gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata", "hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö", "träsk", "flygplats", "industriväg", "trappgata", "kärr", "ringvägen"}
}

func fillInvalidWordsForRoads() {
	invalidWordsForRoads = []string{"väg", "Lv", "Länsväg", "länsväg"}
}

func currentIndexNotLast(index int, strings []string) bool {
	return index < len(strings)-1
}
