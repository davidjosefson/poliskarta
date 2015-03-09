package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"poliskarta/externalservices"
	"poliskarta/filterdescription"
	"poliskarta/filtertitle"
	"strconv"

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

/*
TODO:
	1. Refactor: filtermappar
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
*/

func main() {
	m := martini.Classic()

	m.Group("/", func(r martini.Router) {
		r.Get(":place", allEvents)
		// r.Get(":place/(?P<number>10|[1-9])", singleEvent)
	})

	m.Run()
}

func allEvents(res http.ResponseWriter, req *http.Request, params martini.Params) {
	place, placeErr := isPlaceValid(params["place"])
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
		json := callExternalServicesAndCreateJson(place, limit)
		res.Header().Add("Content-type", "application/json; charset=utf-8")

		//**********************************************
		// Detta behövs medans vi köra allt på localhost,
		// Dålig lösning som är osäker, men då kan vi
		// i alla fall testa allt enkelt
		//**********************************************
		res.Header().Add("Access-Control-Allow-Origin", "*")

		res.Write([]byte(json))
	}
}

func singleEvent(params martini.Params) string {
	return params["number"]
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

func isPlaceValid(parameter string) (string, error) {

	for place, _ := range places {
		if place == parameter {
			return place, nil
		}
	}
	return "", errors.New(parameter + " is not a valid place")

}

func callExternalServicesAndCreateJson(place string, limit int) string {
	policeEvents := externalservices.CallPoliceRSS(places[place], limit)
	filterOutLocationsWords(&policeEvents)
	filterOutTime(&policeEvents)
	filterOutEventType(&policeEvents)
	externalservices.CallMapQuest(&policeEvents)
	policeEventsAsJson := encodePoliceEventsToJSON(policeEvents)

	return string(policeEventsAsJson)
}

func filterOutTime(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].Time = filtertitle.GetTime(event.Title)
	}

	*policeEvents = eventsCopy
}

func filterOutEventType(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, event := range eventsCopy.Events {
		eventsCopy.Events[index].EventType = filtertitle.GetEventType(event.Title)
	}

	*policeEvents = eventsCopy
}
func encodePoliceEventsToJSON(policeEvents externalservices.PoliceEvents) []byte {
	policeEventsAsJson, _ := json.Marshal(policeEvents)

	return policeEventsAsJson
}

func filterOutLocationsWords(policeEvents *externalservices.PoliceEvents) {
	eventsCopy := *policeEvents

	for index, _ := range eventsCopy.Events {
		titleWords, err := filtertitle.FilterTitleWords(eventsCopy.Events[index].Title)

		if err != nil {
			eventsCopy.Events[index].HasPossibleLocation = false
		} else {
			eventsCopy.Events[index].HasPossibleLocation = true
			descriptionWords := filterdescription.FilterDescriptionWords(eventsCopy.Events[index].Description)
			removeDuplicatesAndCombinePossibleLocationWords(titleWords, descriptionWords, &eventsCopy.Events[index].PossibleLocationWords)
		}

	}

	*policeEvents = eventsCopy
}

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
