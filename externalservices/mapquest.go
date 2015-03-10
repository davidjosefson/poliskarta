package externalservices

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func CallMapQuest(policeEvents *PoliceEvents) {
	eventsCopy := *policeEvents
	singleQueryMapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=***REMOVED***&outFormat=xml&location="

	//Skapar en waitgroup som i slutet väntar tills alla goroutines är klara
	var wg sync.WaitGroup

	for index, event := range eventsCopy.Events {
		if len(event.LocationWords) > 0 {

			//increments antalet goroutines den ska vänta på
			wg.Add(1)

			//skicka in wg till varje goroutine
			go singleCallGeoLocationService(singleQueryMapURL, &eventsCopy.Events[index], &wg)
		}

	}

	//vänta tills alla är klara
	wg.Wait()
	*policeEvents = eventsCopy
}

func singleCallGeoLocationService(mapURL string, policeEvent *PoliceEvent, wg *sync.WaitGroup) {
	eventCopy := *policeEvent
	for i := 0; i < len(eventCopy.LocationWords); i++ {
		wordsToSearchWith := URLifyString(eventCopy.LocationWords[i:])

		httpResponse, _ := http.Get(mapURL + wordsToSearchWith)
		xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)

		defer httpResponse.Body.Close()

		geoLocation := geolocationXMLtoStructs(xmlResponse)

		resultIsGood := evaluateGeoLocation(geoLocation)
		// fmt.Println("Searching with: ", wordsToSearchWith)
		if resultIsGood {
			eventCopy.Latitude = geoLocation.Locations[0].LocationAlternatives[0].Latitude
			eventCopy.Longitude = geoLocation.Locations[0].LocationAlternatives[0].Longitude
			eventCopy.CoordinateSearchWords = eventCopy.LocationWords[i:]
			// fmt.Println("Results are good: ", geoLocation.Locations)
			break
		} else {
			// fmt.Println("Results are bad: ", geoLocation.Locations)
		}

	}

	*policeEvent = eventCopy
	defer wg.Done()
}

func URLifyString(sliceToURLify []string) string {
	str := ""
	for _, word := range sliceToURLify {
		str += word + " "
	}
	str = url.QueryEscape(str)
	str = strings.TrimSuffix(str, "+")

	return str
}

func evaluateGeoLocation(geoLocation GeoLocation) bool {
	if geoLocation.Locations[0].LocationAlternatives != nil {
		return true
	} else {
		return false
	}
}

func geolocationXMLtoStructs(XMLresponse []byte) GeoLocation {
	var geoLocation GeoLocation
	xml.Unmarshal(XMLresponse, &geoLocation)

	return geoLocation
}

type GeoLocation struct {
	Locations []Location `xml:"results>result"`
	// ThumbMaps string     `xml:"options>thumbMaps"`
}

type Location struct {
	LocationAlternatives []LocationAlternative `xml:"locations>location"`
}

type LocationAlternative struct {
	Quality   string  `xml:"geocodeQuality"`
	Latitude  float32 `xml:"displayLatLng>latLng>lat"`
	Longitude float32 `xml:"displayLatLng>latLng>lng"`
}
