package externalservices

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func CallPoliceRSS(url string, numEvents int) PoliceEvents {
	httpResponse, _ := http.Get(url)
	xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)

	defer httpResponse.Body.Close()

	policeEvents := policeXMLtoStructs(xmlResponse)

	limitNumOfPoliceEvents(&policeEvents, numEvents)

	return policeEvents
}

func policeXMLtoStructs(policeRSSxml []byte) PoliceEvents {
	var policeEvents PoliceEvents
	xml.Unmarshal(policeRSSxml, &policeEvents)

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
