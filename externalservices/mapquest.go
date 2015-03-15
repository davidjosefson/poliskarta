package externalservices

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func CallMapQuest(policeEvent *PoliceEvent, wg *sync.WaitGroup) {
	// eventCopy := *policeEvent
	mapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=***REMOVED***&outFormat=xml&location="

	if len(policeEvent.LocationWords) > 0 {
		for i := 0; i < len(policeEvent.LocationWords); i++ {
			wordsToSearchWith := URLifyString(policeEvent.LocationWords[i:])

			//***************************
			//
			//		Felhantering behövs
			//
			//***************************

			httpResponse, _ := http.Get(mapURL + wordsToSearchWith)
			xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)

			defer httpResponse.Body.Close()

			geoLocation := geolocationXMLtoStructs(xmlResponse)

			resultIsGood := evaluateGeoLocation(geoLocation)
			// fmt.Println("Searching with: ", wordsToSearchWith)
			if resultIsGood {
				policeEvent.Latitude = geoLocation.Locations[0].LocationAlternatives[0].Latitude
				policeEvent.Longitude = geoLocation.Locations[0].LocationAlternatives[0].Longitude
				policeEvent.CoordinateSearchWords = policeEvent.LocationWords[i:]
				// fmt.Println("Results are good: ", geoLocation.Locations)
				break
			} else {
				// eventCopy.LocationWords = append(eventCopy.LocationWords, "BAD RESULTS")
			}

		}
	}

	// *policeEvent.LocationWords = append(*policeEvent.LocationWords, "FICK INGA KOORD: MapQ")
	// *policeEvent = eventCopy
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
