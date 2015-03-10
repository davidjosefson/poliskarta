package filter

import "poliskarta/externalservices"

func FilterPoliceEvents(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents
	filterOutTime(&eventsCopy)
	filterOutEventType(&eventsCopy)
	filterOutLocationsWords(&eventsCopy)
	*policeEvents = eventsCopy
}

func filterOutTime(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].Time = GetTime(event.Title)
	}

	*policeEvents = eventsCopy
}

func filterOutEventType(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].EventType = GetEventType(event.Title)
	}

	*policeEvents = eventsCopy
}

func filterOutLocationsWords(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, _ := range eventsCopy.Events {
		titleWords, err := FilterTitleWords(eventsCopy.Events[index].Title)

		if err != nil {
			eventsCopy.Events[index].HasPossibleLocation = false
		} else {
			eventsCopy.Events[index].HasPossibleLocation = true
			descriptionWords := FilterDescriptionWords(eventsCopy.Events[index].Description)
			removeDuplicatesAndCombinePossibleLocationWords(titleWords, descriptionWords, &eventsCopy.Events[index].PossibleLocationWords)
		}

	}

	*policeEvents = eventsCopy
}

func removeDuplicatesAndCombinePossibleLocationWords(titleWords []string, descriptionWords []string, locationWords *[]string) {
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
