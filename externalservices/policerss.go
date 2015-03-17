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
	defer httpResponse.Body.Close()

	//If we get http-error when calling the police-RSS
	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)

	//If we get error while reading the response body
	if ioErr != nil {
		return PoliceEvents{}, ioErr
	}

	policeEvents := policeXMLtoStructs(xmlResponse)

	//If the rss-URL is faulty, we will get a 200 OK response, and the only way to know if it IS faulty is
	//that the policeEvents-struct is empty
	if len(policeEvents.Events) < 1 {
		return PoliceEvents{}, errors.New("Communication error with polisen.se (might be a faulty rss-URL)")
	}

	limitNumOfPoliceEvents(&policeEvents, numEvents)

	addAreaToEvents(area, &policeEvents)
	addEventURIs(&policeEvents)

	var err error
	return policeEvents, err
}

//Returns a PoliceEvents instead of PoliceEvent because we want to be able to reuse filter functions
//which only accepts PoliceEvents
func CallPoliceRSSGetSingle(url string, area string, eventID uint32) (PoliceEvents, error) {
	httpResponse, httpErr := http.Get(url)
	defer httpResponse.Body.Close()

	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)
	if ioErr != nil {
		return PoliceEvents{}, ioErr
	}

	//Get police events
	policeEvents := policeXMLtoStructs(xmlResponse)

	//If the rss-URL is faulty, we will get a 200 OK response, and the only way to know if it IS faulty is
	//that the policeEvents-struct is empty
	if len(policeEvents.Events) < 1 {
		return PoliceEvents{}, errors.New("Communication error with polisen.se (might be a faulty rss-URL)")
	}

	//Check if eventID is found among the events
	eventsSingle, idNotFoundErr := findEvent(eventID, policeEvents)

	//Add area-value to event
	addAreaToEvents(area, &eventsSingle)

	addEventURIs(&eventsSingle)

	return eventsSingle, idNotFoundErr
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

	return PoliceEvents{}, &IdNotFoundError{fmt.Sprintf("%d didn't match any events", eventID)}
}

func addAreaToEvents(area string, policeEvents *PoliceEvents) {
	for i, _ := range policeEvents.Events {
		policeEvents.Events[i].AreaValue = area
	}
}
