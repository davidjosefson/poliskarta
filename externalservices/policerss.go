package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"poliskarta/structs"
)

func CallPoliceRSSGetAll(area structs.Area, numEvents int) (structs.PoliceEvents, error) {
	httpResponse, httpErr := http.Get(area.RssURL)
	defer httpResponse.Body.Close()

	//If we get http-error when calling the police-RSS
	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)

	//If we get error while reading the response body
	if ioErr != nil {
		return structs.PoliceEvents{}, ioErr
	}

	policeEvents := policeXMLtoStructs(xmlResponse)

	//If the rss-URL is faulty, we will get a 200 OK response, and the only way to know if it IS faulty is
	//that the policeEvents-struct is empty
	if len(policeEvents.Events) < 1 {
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se (might be a faulty rss-URL)")
	}

	limitNumOfPoliceEvents(&policeEvents, numEvents)

	addAreaToEvents(area, &policeEvents)
	addEventURIs(&policeEvents)

	var err error
	return policeEvents, err
}

//Returns a PoliceEvents instead of PoliceEvent because we want to be able to reuse filter functions
//which only accepts PoliceEvents
func CallPoliceRSSGetSingle(area structs.Area, eventID uint32) (structs.PoliceEvents, error) {
	httpResponse, httpErr := http.Get(area.RssURL)
	defer httpResponse.Body.Close()

	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	xmlResponse, ioErr := ioutil.ReadAll(httpResponse.Body)
	if ioErr != nil {
		return structs.PoliceEvents{}, ioErr
	}

	//Get police events
	policeEvents := policeXMLtoStructs(xmlResponse)

	//If the rss-URL is faulty, we will get a 200 OK response, and the only way to know if it IS faulty is
	//that the policeEvents-struct is empty
	if len(policeEvents.Events) < 1 {
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se (might be a faulty rss-URL)")
	}

	//Check if eventID is found among the events
	eventsSingle, idNotFoundErr := findEvent(eventID, policeEvents)

	//Add area-value to event
	addAreaToEvents(area, &eventsSingle)

	addEventURIs(&eventsSingle)

	return eventsSingle, idNotFoundErr
}

func policeXMLtoStructs(policeRSSxml []byte) structs.PoliceEvents {
	var policeEvents structs.PoliceEvents
	xml.Unmarshal(policeRSSxml, &policeEvents)
	addHashAsID(&policeEvents)

	return policeEvents
}

func limitNumOfPoliceEvents(policeEvents *structs.PoliceEvents, numEvents int) {
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

func addHashAsID(policeEvents *structs.PoliceEvents) {
	eventsCopy := *policeEvents
	for index, _ := range eventsCopy.Events {
		hash := fnv.New32()
		hash.Write([]byte(eventsCopy.Events[index].PoliceEventURL))
		eventsCopy.Events[index].ID = hash.Sum32()
	}
	*policeEvents = eventsCopy
}

func addEventURIs(policeEvents *structs.PoliceEvents) {
	for index, event := range policeEvents.Events {
		link := structs.Link{"self", fmt.Sprintf("http://localhost:3000/areas/%v/%d", event.Area.Value, event.ID)}
		policeEvents.Events[index].Links = append(policeEvents.Events[index].Links, link)
	}
}

func findEvent(eventID uint32, policeEvents structs.PoliceEvents) (structs.PoliceEvents, error) {
	var err error
	for _, event := range policeEvents.Events {
		if eventID == event.ID {
			events := structs.PoliceEvents{}
			events.Events = append(events.Events, event)
			return events, err
		}
	}

	return structs.PoliceEvents{}, &structs.IdNotFoundError{fmt.Sprintf("%d didn't match any events", eventID)}
}

func addAreaToEvents(area structs.Area, policeEvents *structs.PoliceEvents) {
	for i, _ := range policeEvents.Events {
		policeEvents.Events[i].Area = structs.PoliceEventArea{area.Name, area.Value, area.Links}
	}
}
