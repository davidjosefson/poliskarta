package externalservices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func CallMapQuest(policeEvent *PoliceEvent, wg *sync.WaitGroup) {
	// eventCopy := *policeEvent
	mapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=***REMOVED***&outFormat=xml&location="
	defer wg.Done()

	if len(policeEvent.LocationWords) > 0 {
		for i := 0; i < len(policeEvent.LocationWords); i++ {
			wordsToSearchWith := URLifyString(policeEvent.LocationWords[i:])

			httpResponse, httpErr := http.Get(mapURL + wordsToSearchWith)
			defer httpResponse.Body.Close()

			var xmlResponse []byte
			var ioErr error

			if httpErr != nil {
				policeEvent.Latitude = 0
				policeEvent.Longitude = 0
				policeEvent.CoordinateSearchWords = append(policeEvent.CoordinateSearchWords, "<N/A>")
				return
			} else {
				xmlResponse, ioErr = ioutil.ReadAll(httpResponse.Body)

				if ioErr != nil {
					policeEvent.Latitude = 0
					policeEvent.Longitude = 0
					policeEvent.CoordinateSearchWords = append(policeEvent.CoordinateSearchWords, "<N/A>")
					break
				} else {
					geoLocation := geolocationXMLtoStructs(xmlResponse)

					fmt.Println("Geolocation: ", geoLocation)

					resultIsGood, connectErr := evaluateGeoLocation(geoLocation)

					if connectErr != nil {
						policeEvent.Latitude = 0
						policeEvent.Longitude = 0
						policeEvent.CoordinateSearchWords = append(policeEvent.CoordinateSearchWords, "<N/A>")
						break
					} else if resultIsGood {
						policeEvent.Latitude = geoLocation.Locations[0].LocationAlternatives[0].Latitude
						policeEvent.Longitude = geoLocation.Locations[0].LocationAlternatives[0].Longitude
						policeEvent.CoordinateSearchWords = policeEvent.LocationWords[i:]
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

func evaluateGeoLocation(geoLocation GeoLocation) (bool, error) {
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

func geolocationXMLtoStructs(XMLresponse []byte) GeoLocation {
	var geoLocation GeoLocation
	err := xml.Unmarshal(XMLresponse, &geoLocation)

	if err != nil {
		fmt.Println("Geo XML-Struct-error: ", err.Error())
	}

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
