package main

import (
	"fmt"
	"poliskarta/helperfunctions"
	"strconv"
	"strings"
)

var EuorpeRoads []string
var arrayOfTitles []string
var namesForPlacesToInclude []string

func main() {
	fillEuropeRoads()
	fillArrayOfTitles()
	fillNamesForPlacesToInclude()

	for _, title := range arrayOfTitles {

		fmt.Println()
		fmt.Println("-------------------------")

		fmt.Print("Processing title: ")
		fmt.Println(title)
		result := rule1(title)
		fmt.Println()
		fmt.Print("Results: ", result)

	}

}

//Rule 1:
func rule1(title string) string {
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
				fmt.Println(placeDescription)
			}

		}
	}

	return "resultat"
}
func fillEuropeRoads() {
	EuorpeRoads = []string{"E4", "E6", "E10", "E12", "E14", "E16", "E18", "E22", "E45", "E65", "E", "Lv"}
}
func fillArrayOfTitles() {
	arrayOfTitles = []string{
		"Rökutveckling, Blåregnsgatan.",
		"Inbrott i bostad, Sösdala.",
		"Personbilar i kollision på E6.",
		"Bilbrand på Tians väg",
		"Butiksrån, Norra Grängesbergsgatan.",
		"Flera fordon i kollision på E22, Stavröd.",
		"Två män gripna misstänkta för stöld, Hamnen Trelleborg",
		"Inbrott i villa, Vellinge-Månstorp.",
		"Rån mot apotek, Gösta Lundhs gata.",
		"Lastbil välter, Inre ringvägen, i höjd med Sege industriområde.", //Hitta ord som innehåller "vägen"
		"Två bilar krockar, korsningen Klörupsvägen och Havrejordsvägen.",
		"Två personbilar i sidokollision, E22, Bäckaskog.",
		"Personrån.",
		"Stopp och kontroll av stulen bil, väg 111, Laröd.",
		"Trafikolycka/rattfylleri, Vankivavägen, Hässleholm.",
		"Trafifkolycka, Länsväg 769, mellan Skurup och Rydsgård.",
		"Trafikolycka på Länsväg 1329, Norra Rörum.",
		"Två bilar krockar, korsningen Skolgatan och Andra Avenyn.",
		"Person påkörd, Storgatan. Senare ändrad till sjukdom/olycksfall.",
		"Inbrott i företagslokal, Fosie industriområde.", //Ta bara fosie
		"Två personbilar i kollision, Hjalmar Brantings väg.",
		"Brand i sopstation på Ramels väg i Malmö",
		"Misshandel, Anna Lindhs plats.",
		"Misshandel på Onsala danska väg i Gottskär. Detta var inte bra.",
		"Något hände på Lv 985, Hakarp, Huskvarna.",
		"Georg Lückligs väg och larm om ungdomar som bråkar. På platsen.",
		"personbil-minibuss på vägen mellan Hinneryd-Traryd. I minibussen",
		"Per Högströmsgatan",
		"Singelolycka, väg 1560, Tomelillabygden.",
		"Skit hände på Trumvägen område Hamptjärnsmoran Boden är en åra",
	}
}

func fillNamesForPlacesToInclude() {
	namesForPlacesToInclude = []string{"väg", "gränd", "plats", "gata", "led", "torg", "park", "trappa", "trappor", "bro", "gångbro", "allé", "alle", "aveny", "plan", "kaj", "hamn", "strand", "stig", "backe", "kajen", "hamnen", "holme", "holmar", "dockan", "parkväg", "byväg", "byaväg", "gård", "stråket", "tvärgata", "gårdar", "parkgata", "idrottsväg", "broväg", "vägen", "stationsgata", "hamngata", "bangårdsgata", "fätåg", "kyrkogata", "hage", "stråket", "ö", "träsk", "flygplats", "industriväg", "trappgata", "kärr"}
}
