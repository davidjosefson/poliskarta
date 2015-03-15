package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
)

func CallPoliceRSSGetAll(url string, area string, numEvents int) (PoliceEvents, error) {
	httpResponse, httpErr := http.Get(url)
	if httpErr != nil {
		return PoliceEvents{}, httpErr
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)
	if ioErr != nil {
		return PoliceEvents{}, ioErr
	}
	defer httpResponse.Body.Close()

	policeEvents := policeXMLtoStructs(xmlResponse)

	limitNumOfPoliceEvents(&policeEvents, numEvents)

	addAreaToEvents(area, &policeEvents)
	addEventURIs(&policeEvents)

	return policeEvents, error{}
}

//Returns a PoliceEvents instead of PoliceEvent because we want to be able to reuse filter functions
//which only accepts PoliceEvents
func CallPoliceRSSGetSingle(url string, area string, eventID uint32) (PoliceEvents, error) {
	httpResponse, httpErr := http.Get(url)
	if httpErr != nil {
		return PoliceEvents{}, httpErr
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)
	if ioErr != nil {
		return PoliceEvents{}, ioErr
	}
	defer httpResponse.Body.Close()

	//Get police events
	policeEvents := policeXMLtoStructs(xmlResponse)

	//Check if eventID is found among the events
	eventsSingle, err := findEvent(eventID, policeEvents)

	//Add area-value to event
	addAreaToEvents(area, &eventsSingle)

	addEventURIs(&eventsSingle)

	return eventsSingle, err
}

func policeXMLtoStructs(policeRSSxml []byte) PoliceEvents {
	var policeEvents PoliceEvents
	xml.Unmarshal(policeRSSxml, &policeEvents)
	addHashAsID(&policeEvents)

	return policeEvents
}

func limitNumOfPoliceEvents(policeEvents *PoliceEvents, numEvents int) {
	copyEvents := *policeEvents

	//Limit maximum num of events to 50
	if numEvents > 50 {
		numEvents = 50
	}

	//Limit number of events to requested amount
	if numEvents < len(policeEvents.Events) {
		copyEvents.Events = copyEvents.Events[:numEvents]
	}

	*policeEvents = copyEvents
}

func addHashAsID(policeEvents *PoliceEvents) {
	eventsCopy := *policeEvents
	for index, _ := range eventsCopy.Events {
		hash := fnv.New32()
		hash.Write([]byte(eventsCopy.Events[index].Link))
		eventsCopy.Events[index].ID = hash.Sum32()
	}
	*policeEvents = eventsCopy
}

func addEventURIs(policeEvents *PoliceEvents) {
	for index, event := range policeEvents.Events {
		eventURI := fmt.Sprintf("http://localhost:3000/areas/%v/%d", event.AreaValue, event.ID)
		policeEvents.Events[index].EventURI = eventURI
	}
}

func findEvent(eventID uint32, policeEvents PoliceEvents) (PoliceEvents, error) {
	var err error
	for _, event := range policeEvents.Events {
		if eventID == event.ID {
			events := PoliceEvents{}
			events.Events = append(events.Events, event)
			return events, err
		}
	}

	err = errors.New(string(eventID) + " didn't match any events")
	return PoliceEvents{}, err
}

func addAreaToEvents(area string, policeEvents *PoliceEvents) {
	for i, _ := range policeEvents.Events {
		policeEvents.Events[i].AreaValue = area
	}
}
