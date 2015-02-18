package main

import (
	"fmt"
	"strings"
)

func main() {

	arrayOfTitles := []string{"2015-02-15 11:08, Inbrott, Östra Göinge", "2015-01-20 10:28, Trafikolycka, Landskrona", "2015-02-15 00:49, Sammanfattning kväll, Norrbotten", "2015-02-15 00:49, samManfattning kväll, Norrbotten"}

	for _, title := range arrayOfTitles {

		fmt.Println()
		fmt.Println("-------------------------")

		fmt.Print("Processing title: ")
		fmt.Println(title)
		fmt.Println()
		fmt.Print("Results: ")

		hasSamanfattning := checkForSammanfattning(title)

		if hasSamanfattning {
			fmt.Println("Titlen har sammanfattning")

		} else {

			titleWithoutTimeStamp := removeTimeStampFromTitle(title)
			titleLocationInfo := removeNonLocationInfoFromTitle(titleWithoutTimeStamp)
			fmt.Println(titleLocationInfo[0])
		}

	}

}

//Title rule 1: Remove timestamp-thingy
func removeTimeStampFromTitle(title string) string {

	titleWithoutTimeStamp := strings.SplitN(title, ",", 2)

	return strings.TrimSpace(titleWithoutTimeStamp[1])

}

//Title rule 2: Extract location info
func removeNonLocationInfoFromTitle(title string) []string {

	titleWords := strings.Split(title, ",")

	titleLocation := make([]string, 0)

	for i := 1; i < len(titleWords); i++ {

		if startsWithUppercase(titleWords[i]) {
			for j := i; j < len(titleWords); j++ {
				locationToAdd := strings.TrimSpace(titleWords[j])
				titleLocation = append(titleLocation, locationToAdd)
			}
			break
		}

	}

	return titleLocation

}

//Title Rule 3: Dont check for location if sammanfattning exists in title
func checkForSammanfattning(title string) bool {
	return strings.Contains(strings.ToLower(title), "sammanfattning")
}

//Helper func
func startsWithUppercase(s string) bool {
	return (string(s[0]) == strings.ToUpper(string(s[0])))
}

/*func main() {

	handler("2015-02-18 02:09, Brand, Burlöv.", "Bilbrand på Parkgatan, Arlöv.")
	handler("2015-02-15 11:08, Inbrott, Östra Göinge", "Inbrott i fritidsbostad, Knislinge.")
	handler("2015-02-16 13:55, Trafikolycka, Kristianstad", "Personbilar i kollision på Blekingevägen.")
	handler("2015-02-17 23:32, Rattfylleri, Landskrona.", "Rattfylleri, olovlig körning, Pumpgatan, Landskrona.")
	handler("2015-02-11 18:33, Trafikhinder, Kävlinge", "Tankbil blockerar del av vägbana, E6, Hofterup.")
	handler("2015-01-20 10:28, Trafikolycka, Landskrona", "Två personbilar i kollision, Hjalmar Brantings väg.")

}
func handler(title string, desc string) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println(desc)
	//Get and trim last word from title
	splitTitle := strings.Split(title, " ")
	title2ndLastWord := strings.TrimSpace(splitTitle[len(splitTitle)-2])
	titleLastWord := strings.TrimSpace(splitTitle[len(splitTitle)-1])
	titleLastWord = strings.TrimSuffix(titleLastWord, ".")

	//Get and trim last and 2nd to last word from description
	splitDesc := strings.Split(desc, " ")
	desc2ndLastword := strings.TrimSpace(splitDesc[len(splitDesc)-2])
	desc2ndLastword = strings.TrimSuffix(desc2ndLastword, ",")
	descLastword := strings.TrimSpace(splitDesc[len(splitDesc)-1])
	descLastword = strings.TrimSuffix(descLastword, ".")

	address := ""

	fmt.Println()
	fmt.Println("Address is: ")

	//determine if 1 or 2 words from title describes the location
	if title2ndLastWord == strings.TrimSuffix(title2ndLastWord, ",") {
		address = title2ndLastWord + " " + titleLastWord
	} else {
		address = titleLastWord
	}
	/*
		//determine if description repeats location from title, and if 1 or 2 words from desc describes the location
		if address != descLastword && startsWithUppercase(descLastword) {
			address = descLastword + " " + address
		}

		if startsWithUppercase(desc2ndLastword) {
			address = desc2ndLastword + " " + address
		}



	//determine if description repeats location from title, and if 1 or 2 words from desc describes the location
	if address != descLastword {
		if startsWithUppercase(descLastword) {
			address = descLastword + " " + address
		}

	} else if startsWithUppercase(desc2ndLastword) {
		address = desc2ndLastword + " " + descLastword + address
	}

	fmt.Println(address)
	fmt.Println()
	fmt.Println("----------------------------------")

}

func startsWithUppercase(s string) bool {
	return (string(s[0]) == strings.ToUpper(string(s[0])))
}
*/
