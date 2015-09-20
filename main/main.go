package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"poliskarta/externalservices"
	"poliskarta/filter"
	"poliskarta/structs"
	"strconv"
	"sync"

	"github.com/go-martini/martini"
)

var areas = structs.AreasStruct{areasArray, []structs.Link{structs.Link{"self", structs.APIURL + "areas"}}}

var areasArray = []structs.Area{
	structs.Area{"Blekinge", "blekinge", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss", 56.283333, 15.116667, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/blekinge"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss"}}},
	structs.Area{"Dalarna", "dalarna", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Dalarna/?feed=rss", 60.678611, 15.600556, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/dalarna"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Dalarna/?feed=rss"}}},
	structs.Area{"Gotland", "gotland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gotland/?feed=rss", 57.499167, 18.509444, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/gotland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gotland/?feed=rss"}}},
	structs.Area{"Gävleborg", "gavleborg", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gavleborg/?feed=rss", 60.780556, 16.655278, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/gavleborg"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gavleborg/?feed=rss"}}},
	structs.Area{"Halland", "halland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Halland/?feed=rss", 56.716667, 12.821111, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/halland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Halland/?feed=rss"}}},
	structs.Area{"Jämtland", "jamtland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jamtland/?feed=rss", 63.283056, 14.238333, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/jamtland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jamtland/?feed=rss"}}},
	structs.Area{"Jönköping", "jonkoping", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jonkoping/?feed=rss", 57.750000, 14.200000, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/jonkoping"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Jonkoping/?feed=rss"}}},
	structs.Area{"Kalmar", "kalmar", "https://polisen.se/Kalmar_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Kalmar-lan/?feed=rss", 56.733333, 15.9, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/kalmar"}, structs.Link{"origin", "https://polisen.se/Kalmar_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Kalmar-lan/?feed=rss"}}},
	structs.Area{"Kronoberg", "kronoberg", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Kronoberg?feed=rss", 56.79, 14.44, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/kronoberg"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Kronoberg?feed=rss"}}},
	structs.Area{"Norrbotten", "norrbotten", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Norrbotten?feed=rss", 67.135833, 18.501111, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/norrbotten"}, structs.Link{"origin", "ttps://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Norrbotten?feed=rss"}}},
	structs.Area{"Skåne", "skane", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Skane?feed=rss", 56, 13.45, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/skane"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Skane?feed=rss"}}},
	structs.Area{"Stockholm", "stockholm", "https://polisen.se/Stockholms_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Stockholms-lan/?feed=rss", 59.333333, 18.166667, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/stockholm"}, structs.Link{"origin", "https://polisen.se/Stockholms_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Stockholms-lan/?feed=rss"}}},
	structs.Area{"Södermanland", "sodermanland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Sodermanland?feed=rss", 58.771111, 16.869444, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/sodermanland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Sodermanland?feed=rss"}}},
	structs.Area{"Uppsala", "uppsala", "https://polisen.se/Uppsala_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Uppsala-lan/?feed=rss", 59.858333, 17.65, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/uppsala"}, structs.Link{"origin", "https://polisen.se/Uppsala_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Uppsala-lan/?feed=rss"}}},
	structs.Area{"Värmland", "varmland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Varmland?feed=rss", 59.425556, 13.271389, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/varmland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Varmland?feed=rss"}}},
	structs.Area{"Västerbotten", "vasterbotten", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasterbotten?feed=rss", 64.344722, 18.314167, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/vasterbotten"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasterbotten?feed=rss"}}},
	structs.Area{"Västernorrland", "vasternorrland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasternorrland?feed=rss", 62.733333, 16.933333, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/vasternorrland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vasternorrland?feed=rss"}}},
	structs.Area{"Västmanland", "vastmanland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastmanland?feed=rss", 59.645556, 16.424444, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/vastmanland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastmanland?feed=rss"}}},
	structs.Area{"Västra Götaland", "vastragotaland", "https://polisen.se/Vastra_Gotaland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastra-Gotaland/?feed=rss", 58.216944, 11.733333, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/vastragotaland"}, structs.Link{"origin", "https://polisen.se/Vastra_Gotaland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastra-Gotaland/?feed=rss"}}},
	structs.Area{"Örebro", "orebro", "https://polisen.se/Orebro_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Orebro-lan/?feed=rss", 59.266667, 15.216667, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/orebro"}, structs.Link{"origin", "https://polisen.se/Orebro_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Orebro-lan/?feed=rss"}}},
	structs.Area{"Östergötland", "ostergotland", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Ostergotland?feed=rss", 58.410447, 15.613558, 8,
		[]structs.Link{structs.Link{"self", structs.APIURL + "areas/ostergotland"}, structs.Link{"origin", "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Ostergotland?feed=rss"}}},
}

/*
Mjöliga förbättringar:
	1. ta bort "född -90" och andra "-<årtal>"
	2. Kolla om första order innehåller något typ "väg" eller "gatan", för då ska det inte tas bort.
	3. Rensa bort alla /n
	4. Separera interna structs och JSON-structs
	5. PoliceRSS: ändra namn på policeXMLToStructs och lägg in 	AddEvents och AddArea-metoderna till policeXMLToStructs,
	så de inte behöver ligga dubbelt
*/

func main() {
	m := martini.Classic()

	m.Get("/api/v1/areas", allAreas)
	m.Get("/api/v1/areas/:place", allEvents)
	m.Get("/api/v1/areas/:place/:eventid", singleEvent)

	m.Run()
}

func allAreas(res http.ResponseWriter, req *http.Request) {

	//Encodes areas to JSON
	json, _ := json.Marshal(areas)

	res.Header().Add("Content-type", "application/json; charset=utf-8")

	//**********************************************
	// Detta behövs medans vi köra allt på localhost,
	// Dålig lösning som är osäker, men då kan vi
	// i alla fall testa allt enkelt
	//**********************************************
	res.Header().Add("Access-Control-Allow-Origin", "*")

	res.Write(json)
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
		json, connectErr := callPoliceRSSGetJSONAllEvents(area, limit)

		if connectErr != nil {
			status := http.StatusInternalServerError
			res.WriteHeader(status) // http-status 500
			errorMessage := fmt.Sprintf("%v: %v \n\n%v", status, http.StatusText(status), connectErr.Error())
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

func singleEvent(res http.ResponseWriter, req *http.Request, params martini.Params) {
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
		json, connectErr := callPoliceRSSGetJSONSingleEvent(area, uint32(eventID))
		if connectErr != nil {
			if idNotFoundErr, ok := connectErr.(*structs.IdNotFoundError); ok {
				//If id not found - return 400-error
				res.WriteHeader(http.StatusBadRequest) // http-status 400
				errorMessage := fmt.Sprintf("%v: %v \n\n%v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), idNotFoundErr.Error())
				res.Write([]byte(errorMessage))
			} else {
				//If other error, return 500-error
				res.WriteHeader(http.StatusBadRequest) // http-status 400
				errorMessage := fmt.Sprintf("%v: %v \n\n%v", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), connectErr.Error())
				res.Write([]byte(errorMessage))
			}
		} else {
			//If no error, return values
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
	limit := 1000
	var err error
	if param != "" {
		limit, err = strconv.Atoi(param)
		if err != nil {
			err = errors.New(param + " is not a valid positive number")
		}
		if limit < 1 {
			err = errors.New(param + " is not a positive number")
		} else if limit > 50 {
			limit = 500
		}
	}

	return limit, err
}

func isPlaceValid(parameter string) (structs.Area, error) {
	for _, area := range areas.Areas {
		if area.Value == parameter {
			return area, nil
		}
	}

	return structs.Area{}, errors.New(parameter + " is not a valid place")
}

func isEventIDValid(parameter string) (uint64, error) {
	id, err := strconv.ParseUint(parameter, 10, 32)

	if err != nil {
		err = errors.New(parameter + " is not a valid event-ID")
	}

	return id, err
}

func callPoliceRSSGetJSONAllEvents(area structs.Area, limit int) ([]byte, error) {
	policeEvents, err := externalservices.CallPoliceRSSGetAll(area, limit)
	if err != nil {
		return []byte{}, err
	}

	filter.FilterPoliceEvents(&policeEvents)
	policeEventsAsJson := encodePoliceEventsToJSON(policeEvents)

	return policeEventsAsJson, err
}

func callPoliceRSSGetJSONSingleEvent(area structs.Area, eventID uint32) ([]byte, error) {
	policeEvents, err := externalservices.CallPoliceRSSGetSingle(area, eventID)
	if err != nil {
		return []byte{}, err
	}

	//Creating a waitgroup which will wait until all goroutines are finished
	var wg sync.WaitGroup

	//How many goroutines it should wait on
	wg.Add(1)

	go externalservices.CallPoliceScraping(&policeEvents.Events[0], &wg)

	//Is needed before calling MapQuest
	filter.FilterPoliceEvents(&policeEvents)

	//If location-words are present in the event try to find coordinates
	if policeEvents.Events[0].Location != nil {
		wg.Add(1)
		go externalservices.CallMapQuest(&policeEvents.Events[0], &wg)
	}

	//Wait for all goroutines to finish
	wg.Wait()

	policeEventsAsJson := encodePoliceEventToJSON(policeEvents.Events[0])

	return policeEventsAsJson, err
}

func encodePoliceEventsToJSON(policeEvents structs.PoliceEvents) []byte {
	policeEventsAsJson, _ := json.Marshal(policeEvents)

	return policeEventsAsJson
}

func encodePoliceEventToJSON(policeEvent structs.PoliceEvent) []byte {
	policeEventAsJson, _ := json.Marshal(policeEvent)

	return policeEventAsJson
}
