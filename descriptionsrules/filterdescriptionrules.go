package main

import (
	"poliskarta/helperfunctions"
	"strconv"
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
var inValidWordsForRoads []string

//Rule 1:
func Rule1(title string) []string {
	fillEuropeRoads()
	fillValidWordsForPlaces()

	descWords := strings.Split(title, " ")
	trimmedDescWords := helperfunctions.TrimSpacesFromArray(descWords)

	placeDescription := []string{}

	for i := 1; i < len(trimmedDescWords); i++ {

		currentWord := trimmedDescWords[i]

		//ta inte med ordet om förra ordet slutar med "."
		if strings.HasSuffix(trimmedDescWords[i-1], ".") {
			continue
		}

		if (currentWord == "väg" || currentWord == "Lv" || currentWord == "Länsväg") && i < len(trimmedDescWords)-1 {

			wordToAdd := helperfunctions.TrimSuffixFromWord(trimmedDescWords[i+1], ".")
			wordToAdd = helperfunctions.TrimSuffixFromWord(wordToAdd, ",")

			if _, err := strconv.Atoi(wordToAdd); err == nil {

				placeDescription = append(placeDescription, wordToAdd)
				continue
			}
		}

		//Ta eventuellt med ordet om det börjar med stor bokstav
		if helperfunctions.StartsWithUppercase(currentWord) {

			//Om det inte är en Europaväg så ska det eventuellt vara med
			if !helperfunctions.StringInSliceIgnoreCase(currentWord, europeRoads) {

				//ta bort punkter och kommatecken
				currentWord = helperfunctions.TrimSuffixFromWord(currentWord, ".")
				currentWord = helperfunctions.TrimSuffixFromWord(currentWord, ",")

				//Nu får ordet läggas in (tror vi)
				placeDescription = append(placeDescription, currentWord)
			}

		}
	}

	return placeDescription
}
func fillEuropeRoads() {
	europeRoads = []string{"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65", "E", "Lv"}
}

func fillValidWordsForPlaces() {
	validWordsForPlaces = []string{"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa", "trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj", "hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar", "dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata", "gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata", "hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö", "träsk", "flygplats", "industriväg", "trappgata", "kärr"}
}

func fillInvalidWordsForRoads() {
	inValidWordsForRoads = []string{"väg", "Lv", "Länsväg", "länsväg"}
}
