package main

import (
	"poliskarta/helperfunctions"
	"strconv"
	"strings"
)

var EuorpeRoads []string
var arrayOfTitles []string
var namesForPlacesToInclude []string

//Rule 1:
func Rule1(title string) []string {
	fillEuropeRoads()
	fillNamesForPlacesToInclude()

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
			if !helperfunctions.StringInSliceIgnoreCase(currentWord, EuorpeRoads) {

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
	EuorpeRoads = []string{"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65", "E", "Lv"}
}

func fillNamesForPlacesToInclude() {
	namesForPlacesToInclude = []string{"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa", "trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj", "hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar", "dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata", "gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata", "hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö", "träsk", "flygplats", "industriväg", "trappgata", "kärr"}
}
