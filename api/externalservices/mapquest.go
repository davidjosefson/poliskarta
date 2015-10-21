package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"poliskarta/api/structs"
	"strings"
	"sync"
)

func CallMapQuest(policeEvent *structs.PoliceEvent, credentials structs.Credentials, wg *sync.WaitGroup) {
	mapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=" + credentials.Mapquestkey + "&outFormat=xml&location="
	defer wg.Done()

	if len(policeEvent.Location.Words) > 0 {
		for i := 0; i < len(policeEvent.Location.Words); i++ {
			wordsToSearchWith := URLifyString(policeEvent.Location.Words[i:])

			httpResponse, httpErr := http.Get(mapURL + wordsToSearchWith)

			var xmlResponse []byte
			var ioErr error

			if httpErr != nil {
				policeEvent.Location.Latitude = 0
				policeEvent.Location.Longitude = 0
				policeEvent.Location.SearchWords = append(policeEvent.Location.SearchWords, "<N/A>")
				return
			} else {
				defer httpResponse.Body.Close()
				xmlResponse, ioErr = ioutil.ReadAll(httpResponse.Body)

				if ioErr != nil {
					policeEvent.Location.Latitude = 0
					policeEvent.Location.Longitude = 0
					policeEvent.Location.SearchWords = append(policeEvent.Location.SearchWords, "<N/A>")
					break
				} else {
					geoLocation := geolocationXMLtoStructs(xmlResponse)
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
