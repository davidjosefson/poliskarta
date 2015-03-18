package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"poliskarta/structs"
	"strings"
	"sync"
)

func CallMapQuest(policeEvent *structs.PoliceEvent, wg *sync.WaitGroup) {
	// eventCopy := *policeEvent
	mapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=***REMOVED***&outFormat=xml&location="
	defer wg.Done()

	if len(policeEvent.Location.Words) > 0 {
		for i := 0; i < len(policeEvent.Location.Words); i++ {
			wordsToSearchWith := URLifyString(policeEvent.Location.Words[i:])

			httpResponse, httpErr := http.Get(mapURL + wordsToSearchWith)
			defer httpResponse.Body.Close()

			var xmlResponse []byte
			var ioErr error

			if httpErr != nil {
				policeEvent.Location.Latitude = 0
				policeEvent.Location.Longitude = 0
				policeEvent.Location.SearchWords = append(policeEvent.Location.SearchWords, "<N/A>")
				return
			} else {
				xmlResponse, ioErr = ioutil.ReadAll(httpResponse.Body)

				if ioErr != nil {
					policeEvent.Location.Latitude = 0
					policeEvent.Location.Longitude = 0
					policeEvent.Location.SearchWords = append(policeEvent.Location.SearchWords, "<N/A>")
					break
				} else {
					geoLocation := geolocationXMLtoStructs(xmlResponse)

					fmt.Println("Geolocation: ", geoLocation)

					resultIsGood, connectErr := evaluateGeoLocation(geoLocation)

					if connectErr != nil {
						policeEvent.Location.Latitude = 0
						policeEvent.Location.Longitude = 0
						policeEvent.Location.SearchWords = append(policeEvent.Location.SearchWords, "<N/A>")
						break
					} else if resultIsGood {
						policeEvent.Location.Latitude = geoLocation.Locations[0].LocationAlternatives[0].Latitude
						policeEvent.Location.Longitude = geoLocation.Locations[0].LocationAlternatives[0].Longitude
						policeEvent.Location.SearchWords = policeEvent.Location.Words[i:]
						break
					}
				}
			}

		}
	}

	// *policeEvent.LocationWords = append(*policeEvent.LocationWords, "FICK INGA KOORD: MapQ")
	// *policeEvent = eventCopy

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

func evaluateGeoLocation(geoLocation structs.GeoLocation) (bool, error) {
	var err error

	if len(geoLocation.Locations) > 0 {
		if geoLocation.Locations[0].LocationAlternatives != nil {
			return true, err
		} else {
			return false, err
		}
	} else {
		return false, errors.New("Communication error with mapquest.com")
	}
}

func geolocationXMLtoStructs(XMLresponse []byte) structs.GeoLocation {
	var geoLocation structs.GeoLocation
	err := xml.Unmarshal(XMLresponse, &geoLocation)

	if err != nil {
		fmt.Println("Geo XML-Struct-error: ", err.Error())
	}

	return geoLocation
}
