package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"poliskarta/externalservices"
	"poliskarta/filter"
	"strconv"
	"sync"

	"github.com/go-martini/martini"
)

var areas = AreasStruct{areasArray}

type AreasStruct struct {
	Areas []Area `json:"areas"`
}
type Area struct {
	Name            string  `json:"name"`
	Value           string  `json:"value"`
	Url             string  `json:"url"`
	RssURL          string  `json:"-"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	GoogleZoomLevel int     `json:"zoomlevel"`
}

var areasArray = []Area{
	Area{"Blekinge", "blekinge", "/blekinge", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss", 56.283333, 15.116667, 8},
	Area{"Dalarna", "dalarna", "/dalarna", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Dalarna/?feed=rss", 60.678611, 15.600556, 8},
	Area{"Gotland", "gotland", "/gotland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gotland/?feed=rss", 57.499167, 18.509444, 8},
	Area{"Gävleborg", "gavleborg", "/gavleborg", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gavleborg/?feed=rss", 60.780556, 16.655278, 8},
	Area{"Halland", "halland", "/halland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Halland/?feed=rss", 56.716667, 12.821111, 8},
	Area{"Jämtland", "jamtland", "/jamtland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jamtland/?feed=rss", 63.283056, 14.238333, 8},
	Area{"Jönköping", "jonkoping", "/jonkoping", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jonkoping/?feed=rss", 57.750000, 14.200000, 8},
	Area{"Kalmar", "kalmar", "/kalmar", "https://polisen.se/Kalmar_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Kalmar-lan/?feed=rss", 56.733333, 15.9, 8},
	Area{"Kronoberg", "kronoberg", "/kronoberg", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Kronoberg?feed=rss", 56.79, 14.44, 8},
	Area{"Norrbotten", "norrbotten", "/norrbotten", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Norrbotten?feed=rss", 67.135833, 18.501111, 8},
	Area{"Skåne", "skane", "/skane", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Skane?feed=rss", 56, 13.45, 8},
	Area{"Stockholm", "stockholm", "/stockholm", "https://polisen.se/Stockholms_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Stockholms-lan/?feed=rss", 59.333333, 18.166667, 8},
	Area{"Södermanland", "sodermanland", "/sodermanland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Sodermanland?feed=rss", 58.771111, 16.869444, 8},
	Area{"Uppsala", "uppsala", "/uppsala", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Uppsala?feed=rss", 59.858333, 17.65, 8},
	Area{"Värmland", "varmland", "/varmland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Varmland?feed=rss", 59.425556, 13.271389, 8},
	Area{"Västerbotten", "vasterbotten", "/vasterbotten", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasterbotten?feed=rss", 64.344722, 18.314167, 8},
	Area{"Västernorrland", "vasternorrland", "/vasternorrland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasternorrland?feed=rss", 62.733333, 16.933333, 8},
	Area{"Västmanland", "vastmanland", "/vastmanland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastmanland?feed=rss", 59.645556, 16.424444, 8},
	Area{"Västra Götaland", "vastragotaland", "/vastragotaland", "https://polisen.se/Vastra_Gotaland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastra-Gotaland/?feed=rss", 58.216944, 11.733333, 8},
	Area{"Örebro", "orebro", "/orebro", "https://polisen.se/Orebro_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Orebro-lan/?feed=rss", 59.266667, 15.216667, 8},
	Area{"Östergötland", "ostergotland", "/ostergotland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Ostergotland?feed=rss", 58.410447, 15.613558, 8},
}

/*
TODO:
	2. Felhantering: polis/mapquest = nere, errors osv
	3. Stockholms-undantag
	4. Norrbotten: det mesta är fel här! Fler generella regler?
	5. Lägga till HATEOAS på /, där man får en lista över tillgängliga län
		- API-länkar
		- bra namn
		- bra value-namn (utan åäö/mellanslag)
		- koordinater till "mittpunkten"?
	6. Optimera? Antingen:
		- lägga in parameter ?no-coord=true  +  /getCoord/?Norra+gränges
		- databas som sparar tidigare ord+koordinater, så att inget jobb/anrop görs utåt
		- cache i 5 minuter som standard, eller en parameter för att alltid få senaste: ?no-cache=true
	7. omstrukturera policeevents-structen så att den har objekt/grupper av saker
	8. PoliceRSS: ändra namn på policeXMLToStructs och lägg in 	AddEvents och AddArea-metoderna till denna,
	så de inte behöver ligga dubbelt
*/

func main() {
	m := martini.Classic()

	m.Get("/areas/:place", allEvents)
	m.Get("/areas", allAreas)
	m.Get("/areas/:place/:eventid", singleEvent)
	// r.Get(":place/(?P<number>10|[1-9])", singleEvent)

	m.Run()
}

func allAreas(res http.ResponseWriter, req *http.Request) {
	json := encodeAreasToJSON()
	res.Header().Add("Content-type", "application/json; charset=utf-8")

	//**********************************************
	// Detta behövs medans vi köra allt på localhost,
	// Dålig lösning som är osäker, men då kan vi
	// i alla fall testa allt enkelt
	//**********************************************
	res.Header().Add("Access-Control-Allow-Origin", "*")

	res.Write(json)
}

func encodeAreasToJSON() []byte {
	areasAsJSON, err := json.Marshal(areas)
	if err != nil {
		//*********
		//Error som inte hanteras, glöm inte bort.
		//*********

		fmt.Println("encodingerror: ", err.Error())
	}

	return areasAsJSON
}

func allEvents(res http.ResponseWriter, req *http.Request, params martini.Params) {
	area, placeErr := isPlaceValid(params["place"])
	limit, limitErr := isLimitParamValid(req.FormValue("limit"))

	if placeErr != nil {
		status := http.StatusBadRequest
		res.WriteHeader(status) // http-status 400
		errorMessage := fmt.Sprintf("%v: %v \n\n%v", status, http.StatusText(status), placeErr.Error())
		res.Write([]byte(errorMessage))
	} else if limitErr != nil {
		status := http.StatusBadRequest
		res.WriteHeader(status) // http-status 400
		errorMessage := fmt.Sprintf("%v: %v \n\n%v", status, http.StatusText(status), limitErr.Error())
		res.Write([]byte(errorMessage))
	} else {
		json := callPoliceRSSGetJSONAllEvents(area.RssURL, area.Value, limit)
		res.Header().Add("Content-type", "application/json; charset=utf-8")

		//**********************************************
		// Detta behövs medans vi köra allt på localhost,
		// Dålig lösning som är osäker, men då kan vi
		// i alla fall testa allt enkelt
		//**********************************************
		res.Header().Add("Access-Control-Allow-Origin", "*")

		res.Write(json)
	}
}

// Is for: r.Get(":place/(?P<number>10|[1-9])", singleEvent)
// func singleEvent(params martini.Params) string {
// 	return params["number"]
// }

func singleEvent(res http.ResponseWriter, req *http.Request, params martini.Params) {
	/*
		- Kolla så att place är valid
		- Gör http-anrop hos polisen med place
		- Hasha alla polis-urler
		- Jämför med param["eventhash"]
		- Returnera JSON av
	*/
	area, placeErr := isPlaceValid(params["place"])
	eventID, idParseErr := isEventIDValid(params["eventid"])

	if placeErr != nil {
		res.WriteHeader(http.StatusBadRequest) // http-status 400
		errorMessage := fmt.Sprintf("%v: %v \n\n%v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), placeErr.Error())
		res.Write([]byte(errorMessage))
	} else if idParseErr != nil {
		res.WriteHeader(http.StatusBadRequest) // http-status 400
		errorMessage := fmt.Sprintf("%v: %v \n\n%v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), idParseErr.Error())
		res.Write([]byte(errorMessage))
	} else {
		json, idNotFoundErr := callPoliceRSSGetJSONSingleEvent(area.RssURL, area.Value, uint32(eventID))
		if idNotFoundErr != nil {
			res.WriteHeader(http.StatusBadRequest) // http-status 400
			errorMessage := fmt.Sprintf("%v: %v \n\n%v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), idNotFoundErr.Error())
			res.Write([]byte(errorMessage))
		} else {
			res.Header().Add("Content-type", "application/json; charset=utf-8")
			//**********************************************
			// Detta behövs medans vi köra allt på localhost,
			// Dålig lösning som är osäker, men då kan vi
			// i alla fall testa allt enkelt
			//**********************************************
			res.Header().Add("Access-Control-Allow-Origin", "*")
			res.Write(json)
		}
	}

}

func isLimitParamValid(param string) (int, error) {
	limit := 10
	var err error
	if param != "" {
		limit, err = strconv.Atoi(param)
		if err != nil {
			err = errors.New(param + " is not a valid positive number")
		}
		if limit < 1 {
			err = errors.New(param + " is not a positive number")
		} else if limit > 50 {
			limit = 50
		}
	}

	return limit, err
}

func isPlaceValid(parameter string) (Area, error) {
	for _, area := range areas.Areas {
		if area.Value == parameter {
			return area, nil
		}
	}

	return Area{}, errors.New(parameter + " is not a valid place")
}

func isEventIDValid(parameter string) (uint64, error) {
	id, err := strconv.ParseUint(parameter, 10, 32)

	if err != nil {
		err = errors.New(parameter + " is not a valid event-ID")
	}

	return id, err
}

func callPoliceRSSGetJSONAllEvents(url string, area string, limit int) []byte {
	policeEvents := externalservices.CallPoliceRSSGetAll(url, area, limit)
	filter.FilterPoliceEvents(&policeEvents)
	policeEventsAsJson := encodePoliceEventsToJSON(policeEvents)

	return policeEventsAsJson
}
func callPoliceRSSGetJSONSingleEvent(url string, area string, eventID uint32) ([]byte, error) {
	policeEvents, err := externalservices.CallPoliceRSSGetSingle(url, area, eventID)
	if err != nil {
		return []byte{}, err
	}

	//Creating a waitgroup which will wait until all goroutines is finished
	var wg sync.WaitGroup

	//How many goroutines it should wait on
	wg.Add(2)

	go externalservices.CallPoliceScraping(&policeEvents.Events[0], &wg)

	//Is needed before calling MapQuest
	filter.FilterPoliceEvents(&policeEvents)

	go externalservices.CallMapQuest(&policeEvents.Events[0], &wg)

	//Wait for all goroutines to finish
	wg.Wait()

	policeEventsAsJson := encodePoliceEventToJSON(policeEvents.Events[0])

	return policeEventsAsJson, err
}

func encodePoliceEventsToJSON(policeEvents externalservices.PoliceEvents) []byte {
	policeEventsAsJson, _ := json.Marshal(policeEvents)

	return policeEventsAsJson
}

func encodePoliceEventToJSON(policeEvent externalservices.PoliceEvent) []byte {
	policeEventAsJson, _ := json.Marshal(policeEvent)

	return policeEventAsJson
}
