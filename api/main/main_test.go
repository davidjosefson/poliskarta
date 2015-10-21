package main

import (
	"fmt"
	"testing"
)

func Test_removeDuplicatesAndCombineLocationWords(t *testing.T) {
	// (titleWords []string, descriptionWords []string, locationWords *[]string

	titles := []string{"Malmö"}
	description := []string{"Ellstorpsgatan", "Malmö"}
	locationWords := []string{}

	removeDuplicatesAndCombineLocationWords(titles, description, &locationWords)

	fmt.Println(locationWords)

}

func Test_findAndFillLocationWords(t *testing.T) {

	events := []PoliceEvent{
		{Title: "Hej, hopp, Malmö",
			Description: "Det hände skit, rån, Ellstorpsgatan"},
		{Title: "Titel2, Rån, Kronoberg",
			Description: "Desc2, Annat skit hände, Landskrona",
		},
	}
	policeEvents := PoliceEvents{events}

	findAndFillLocationWords(&policeEvents)

	fmt.Println(policeEvents.Events[0].LocationWords)
	fmt.Println(policeEvents.Events[1].LocationWords)

}
