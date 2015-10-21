package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"poliskarta/api/structs"
)

func CallPoliceRSSGetAll(area structs.Area, numEvents int) (structs.PoliceEvents, error) {
	httpResponse, httpErr := http.Get(area.RssURL)

	//If we get http-error when calling the police-RSS
	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	defer httpResponse.Body.Close()

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

	addAreaInfoToResponse(&policeEvents, area)
	addEventLinks(&policeEvents, area)

	var err error
	return policeEvents, err
}

// CallPoliceRSSGetSingle Returns a PoliceEvents instead of PoliceEvent because we want to be able to reuse filter functions
// which only accepts PoliceEvents
func CallPoliceRSSGetSingle(area structs.Area, eventID uint32) (structs.PoliceEvents, error) {
	httpResponse, httpErr := http.Get(area.RssURL)

	if httpErr != nil {
		fmt.Println(httpErr.Error())
		return structs.PoliceEvents{}, errors.New("Communication error with polisen.se")
	}

	defer httpResponse.Body.Close()

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

	//Has to be added because mainfilter checks if this info is for "stockholm", and acts accordingly
	addAreaInfoToResponse(&eventsSingle, area)

	addEventLinks(&eventsSingle, area)

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
		numEvents = 500
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

func addEventLinks(policeEvents *structs.PoliceEvents, area structs.Area) {
	for index, event := range policeEvents.Events {
		selfLink := structs.Link{"self", fmt.Sprintf(structs.APIURL+"areas/%v/%d", area.Value, event.ID)}
		originLink := structs.Link{"origin", event.PoliceEventURL}
		policeEvents.Events[index].Links = append(policeEvents.Events[index].Links, selfLink, originLink)
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
		policeEvents.Events[i].Area = &structs.PoliceEventArea{area.Name, area.Value, area.Links}
	}
}
func addAreaInfoToResponse(policeEvents *structs.PoliceEvents, area structs.Area) {
	policeEvents.Name = area.Name
	policeEvents.Value = area.Value
	policeEvents.Latitude = area.Latitude
	policeEvents.Longitude = area.Longitude
	policeEvents.GoogleZoomLevel = area.GoogleZoomLevel
	policeEvents.Links = area.Links
}
