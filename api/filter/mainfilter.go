package filter

import (
	"poliskarta/api/helperfunctions"
	"poliskarta/api/structs"
	"strings"
)

func FilterPoliceEvents(policeEvents *structs.PoliceEvents) {
	eventsCopy := *policeEvents
	filterOutTime(&eventsCopy)
	filterOutEventType(&eventsCopy)
	if policeEvents.Value == "stockholm" {
		addAllWordsAsLocationWords(&eventsCopy)
	} else {
		filterOutLocationsWords(&eventsCopy)
	}
	*policeEvents = eventsCopy
}

func filterOutTime(policeEvents *structs.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].Time = GetTime(event.Title)
	}

	*policeEvents = eventsCopy
}

func filterOutEventType(policeEvents *structs.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].EventType = GetEventType(event.Title)
	}

	*policeEvents = eventsCopy
}

func filterOutLocationsWords(policeEvents *structs.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, _ := range eventsCopy.Events {
		titleWords, err := FilterTitleWords(eventsCopy.Events[index].Title)

		if err == nil {
			//Must be instantiated before we can append words
			eventsCopy.Events[index].Location = &structs.LocationInfo{}
			descriptionWords := FilterDescriptionWords(eventsCopy.Events[index].DescriptionShort)
			removeDuplicatesAndCombineLocationWords(titleWords, descriptionWords, &eventsCopy.Events[index].Location.Words)
		}
	}

	*policeEvents = eventsCopy
}
func addAllWordsAsLocationWords(policeEvents *structs.PoliceEvents) {
	for index, _ := range policeEvents.Events {
		titleWords, err := FilterTitleWords(policeEvents.Events[index].Title)

		if err == nil {
			//Must be instantiated before we can append words
			policeEvents.Events[index].Location = &structs.LocationInfo{}

			descriptionWords := strings.Split(policeEvents.Events[index].DescriptionShort, ",")
			helperfunctions.TrimSpacesFromArray(&descriptionWords)
			removeDuplicatesAndCombineLocationWords(titleWords, descriptionWords, &policeEvents.Events[index].Location.Words)
		}
	}
}

func removeDuplicatesAndCombineLocationWords(titleWords []string, descriptionWords []string, locationWords *[]string) {
	location := []string{}

	for _, descWord := range descriptionWords {
		location = append(location, descWord)
	}

	wordAlreadyExists := false

	for _, titleWord := range titleWords {
		for _, locationWord := range location {
			if titleWord == locationWord {
				wordAlreadyExists = true
				break
			}
		}
		if !wordAlreadyExists {
			location = append(location, titleWord)
		}
	}

	*locationWords = location
}
