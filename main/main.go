package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"poliskarta/filterdescription"
	"poliskarta/filtertitle"
	"strings"
	"sync"

	"github.com/go-martini/martini"
)

var places = map[string]string{
	"blekinge":       "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss",
	"dalarna":        "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Dalarna/?feed=rss",
	"gotland":        "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gotland/?feed=rss",
	"gavleborg":      "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gavleborg/?feed=rss",
	"halland":        "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Halland/?feed=rss",
	"jamtland":       "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jamtland/?feed=rss",
	"jonkoping":      "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jonkoping/?feed=rss",
	"kalmar":         "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Kalmar?feed=rss",
	"kronoberg":      "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Kronoberg?feed=rss",
	"norrbotten":     "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Norrbotten?feed=rss",
	"skane":          "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Skane?feed=rss",
	"stockholm":      "https://polisen.se/Stockholms_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Stockholms-lan/?feed=rss",
	"sodermanland":   "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Sodermanland?feed=rss",
	"uppsala":        "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Uppsala?feed=rss",
	"varmland":       "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Varmland?feed=rss",
	"vasterbotten":   "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasterbotten?feed=rss",
	"vasternorrland": "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasternorrland?feed=rss",
	"vastmanland":    "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastmanland?feed=rss",
	"vastragotaland": "https://polisen.se/Vastra_Gotaland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastra-Gotaland/?feed=rss",
	"orebro":         "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Orebro?feed=rss",
	"ostergotland":   "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Ostergotland?feed=rss",
}

func main() {
	m := martini.Classic()

	m.Group("/", func(r martini.Router) {
		r.Get(":place", allEvents)
		r.Get(":place/(?P<number>10|[1-9])", singleEvent)
	})

	m.Run()

	// //http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css")))) //To find css-files in the css-folder
	// http.ListenAndServe(":9090", nil)
}

func allEvents(res http.ResponseWriter, req *http.Request, params martini.Params) {
	place := params["place"]

	if isPlaceValid(place) {
		json := callExternalServicesAndCreateJson(place)
		res.Header().Add("Content-type", "application/json; charset=utf-8")

		//**********************************************
		// Detta behövs medans vi köra allt på localhost,
		// Dålig lösning som är osäker, men då kan vi
		// i alla fall testa allt enkelt
		//**********************************************
		res.Header().Add("Access-Control-Allow-Origin", "*")

		res.Write([]byte(json))
	} else {
		status := http.StatusBadRequest
		res.WriteHeader(status) // http-status 400
		errorMessage := fmt.Sprintf("%v: %v \n\n\"%v\" is not a valid place", status, http.StatusText(status), place)
		res.Write([]byte(errorMessage))
	}
}

func singleEvent(params martini.Params) string {
	return params["number"]
}

func isPlaceValid(parameter string) bool {
	for place, _ := range places {
		if place == parameter {
			return true
		}
	}
	return false
}

func callExternalServicesAndCreateJson(place string) string {
	/*
		1. Get Police RSS XML
		2. Save each event as event-struct-array
		3. Fill "searchwords"-fields by using the filters
		4. Get google search results using "searchwords" - save coordinates as fields in struct
		5. Convert search result as JSON and return string
	*/

	policeRSSxml := callPoliceRSS(places[place])
	policeEvents := policeXMLtoStructs(policeRSSxml)

	findAndFillPossibleLocationWords(&policeEvents)
	findAndFillCoordinates(&policeEvents)

	policeEventsAsJson := encodePoliceEventsToJSON(policeEvents)

	return string(policeEventsAsJson)
}

func callPoliceRSS(url string) []byte {
	httpResponse, _ := http.Get(url)
	xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)

	defer httpResponse.Body.Close()

	return xmlResponse
}

func policeXMLtoStructs(policeRSSxml []byte) PoliceEvents {
	var policeEvents PoliceEvents
	xml.Unmarshal(policeRSSxml, &policeEvents)

	return policeEvents
}

func encodePoliceEventsToJSON(policeEvents PoliceEvents) []byte {
	policeEventsAsJson, _ := json.Marshal(policeEvents)

	return policeEventsAsJson
}

/* -- STRUCTS -- */

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	Title                 string `xml:"title"`
	Link                  string `xml:"link"`
	Description           string `xml:"description"`
	HasPossibleLocation   bool
	PossibleLocationWords []string
	HasCoordinates        bool
	CoordinateSearchWords []string
	Longitude             float32
	Latitude              float32
}

func findAndFillPossibleLocationWords(policeEvents *PoliceEvents) {
	eventsCopy := *policeEvents

	for index, _ := range eventsCopy.Events {
		titleWords, err := filtertitle.FilterTitleWords(eventsCopy.Events[index].Title)

		if err != nil {
			eventsCopy.Events[index].HasPossibleLocation = false
		} else {
			eventsCopy.Events[index].HasPossibleLocation = true
			descriptionWords := filterdescription.FilterDescriptionWords(eventsCopy.Events[index].Description)
			removeDuplicatesAndCombinePossibleLocationWords(titleWords, descriptionWords, &eventsCopy.Events[index].PossibleLocationWords)
			//Denna används inte längre
			//AddURLifiedURL(&eventsCopy.Events[index])
		}

	}

	*policeEvents = eventsCopy
}

// func AddURLifiedURL(policeEvent *PoliceEvent) {
// 	eventCopy := *policeEvent
// 	str := ""
// 	for _, word := range eventCopy.LocationWords {
// 		str += word + " "
// 	}
// 	str = url.QueryEscape(str)
// 	str = strings.TrimSuffix(str, "+")

// 	eventCopy.URLifiedLocation = str
// 	*policeEvent = eventCopy
// }

func removeDuplicatesAndCombinePossibleLocationWords(titleWords []string, descriptionWords []string, locationWords *[]string) {
	location := []string{}

	for _, descWord := range descriptionWords {
		location = append(location, descWord)
	}

	wordAlreadyExists := false

	for _, titleWord := range titleWords {
		for _, locationWord := range location {
			if titleWord == locationWord {
				wordAlreadyExists = true
				break
			}
		}
		if !wordAlreadyExists {
			location = append(location, titleWord)
		}
	}

	*locationWords = location
}

/*
TODO:
	1. Joina goroutines (gemensam räknare, channel, eller nåt annat)
	2. Fixa så att varje anrop kan söka om med färre ord, ifall den inte får ett bra resultat
		- om den tex. bara får ut COUNTY, spara resultatet men gör om med färre för att se ifall man kan få
			CITY eller helst STREET
	3. Lägg till koordinaterna i policeEventsen och skicka till klienten
	4. Koppla på GUI och casha in VG.

*/

func findAndFillCoordinates(policeEvents *PoliceEvents) {
	eventsCopy := *policeEvents
	singleQueryMapURL := "http://open.mapquestapi.com/geocoding/v1/address?key=***REMOVED***&outFormat=xml&location="

	//Skapar en waitgroup som i slutet väntar tills alla goroutines är klara
	var wg sync.WaitGroup

	for index, event := range eventsCopy.Events {
		if event.HasPossibleLocation {

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
	eventCopy.HasCoordinates = false
	for i := 0; i < len(eventCopy.PossibleLocationWords); i++ {
		wordsToSearchWith := URLifyString(eventCopy.PossibleLocationWords[i:])

		httpResponse, _ := http.Get(mapURL + wordsToSearchWith)
		xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)

		defer httpResponse.Body.Close()

		geoLocation := geolocationXMLtoStructs(xmlResponse)

		resultIsGood := evaluateGeoLocation(geoLocation)
		fmt.Println("Searching with: ", wordsToSearchWith)
		if resultIsGood {
			eventCopy.HasCoordinates = true
			eventCopy.Latitude = geoLocation.Locations[0].LocationAlternatives[0].Latitude
			eventCopy.Longitude = geoLocation.Locations[0].LocationAlternatives[0].Longitude
			eventCopy.CoordinateSearchWords = eventCopy.PossibleLocationWords[i:]
			fmt.Println("Results are good: ", geoLocation.Locations)
			break
		} else {
			fmt.Println("Results are bad: ", geoLocation.Locations)
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

// func batchCallGeoLocationService(url string) GeoLocationBatch {
// 	httpResponse, _ := http.Get(url)
// 	xmlResponse, _ := ioutil.ReadAll(httpResponse.Body)
// 	// xmlResponse := []byte("<response> <info> <statusCode>0</statusCode> <messages/> <copyright> <imageUrl>http://api.mqcdn.com/res/mqlogo.gif</imageUrl> <imageAltText>© 2015 MapQuest, Inc.</imageAltText> <text>© 2015 MapQuest, Inc.</text> </copyright> </info> <results> <result> <providedLocation> <location>Köping</location> </providedLocation> <locations> <location> <street/> <postalCode/> <geocodeQuality>COUNTY</geocodeQuality> <geocodeQualityCode>A4XXX</geocodeQualityCode> <dragPoint>false</dragPoint> <sideOfStreet>N</sideOfStreet> <displayLatLng> <latLng> <lat>59.579954</lat> <lng>15.879022</lng> </latLng> </displayLatLng> <linkId>0</linkId> <type>s</type> <latLng> <lat>59.579954</lat> <lng>15.879022</lng> </latLng> <mapUrl> <![CDATA[http://open.mapquestapi.com/staticmap/v4/getmap?key=Fmjtd|luu82l6r20,8s=o5-94ralr&type=map&size=225,160&pois=purple-1,59.579954,15.8790217384978,0,0|&center=59.579954,15.8790217384978&zoom=9&rand=582209154 ]]> </mapUrl> </location> <location> <street/> <postalCode/> <geocodeQuality>COUNTY</geocodeQuality> <geocodeQualityCode>A4XXX</geocodeQualityCode> <dragPoint>false</dragPoint> <sideOfStreet>N</sideOfStreet> <displayLatLng> <latLng> <lat>59.513743</lat> <lng>15.997048</lng> </latLng> </displayLatLng> <linkId>0</linkId> <type>s</type> <latLng> <lat>59.513743</lat> <lng>15.997048</lng> </latLng> <mapUrl> <![CDATA[http://open.mapquestapi.com/staticmap/v4/getmap?key=Fmjtd|luu82l6r20,8s=o5-94ralr&type=map&size=225,160&pois=purple-2,59.5137434,15.9970475,0,0|&center=59.5137434,15.9970475&zoom=9&rand=582209154 ]]> </mapUrl> </location> </locations> </result> <result> <providedLocation> <location>Erikslund Västerås</location> </providedLocation> <locations> <location> <street/> <postalCode/> <geocodeQuality>COUNTY</geocodeQuality> <geocodeQualityCode>A4XXX</geocodeQualityCode> <dragPoint>false</dragPoint> <sideOfStreet>N</sideOfStreet> <displayLatLng> <latLng> <lat>59.613284</lat> <lng>16.462255</lng> </latLng> </displayLatLng> <linkId>0</linkId> <type>s</type> <latLng> <lat>59.613284</lat> <lng>16.462255</lng> </latLng> <mapUrl> <![CDATA[http://open.mapquestapi.com/staticmap/v4/getmap?key=Fmjtd|luu82l6r20,8s=o5-94ralr&type=map&size=225,160&pois=purple-1,59.6132838,16.4622554,0,0|&center=59.6132838,16.4622554&zoom=9&rand=582209154 ]]> </mapUrl> </location> </locations> </result> </results> <options> <maxResults>-1</maxResults> <thumbMaps>true</thumbMaps> <ignoreLatLngInput>false</ignoreLatLngInput> <boundingBox/> </options> </response>")
// 	fmt.Println(string(xmlResponse))

// 	defer httpResponse.Body.Close()

// 	geoLocationBatch := geolocationXMLtoStructs(xmlResponse)

// 	return geoLocationBatch
// }

// func geolocationXMLtoStructs(XMLresponse []byte) GeoLocationBatch {
// 	var geoLocationBatch GeoLocationBatch
// 	xml.Unmarshal(XMLresponse, &geoLocationBatch)

// 	return geoLocationBatch
// }

// type GeoLocationBatch struct {
// 	Locations []Location `xml:"results>result"`
// 	ThumbMaps string     `xml:"options>thumbMaps"`
// }

// type Location struct {
// 	LocationAlternatives []LocationAlternative `xml:"locations>location"`
// }

// type LocationAlternative struct {
// 	Quality   string `xml:"geocodeQuality"`
// 	Latitude  string `xml:"displayLatLng>latLng>lat"`
// 	Longitude string `xml:"displayLatLng>latLng>lng"`
// }

/*type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	Title            string `xml:"title"`
	Link             string `xml:"link"`
	Description      string `xml:"description"`
	HasLocation      bool
	LocationWords    []string
	URLifiedLocation string
	Longitude        float32
	Latitude         float32
}

*/
