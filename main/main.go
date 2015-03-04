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
	"vastragotaland": "https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Vastragotaland?feed=rss",
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
	findAndFillLocationWords(&policeEvents)
	// 4. Get google search results using "searchwords" - save coordinates as fields in struct
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
	Title            string `xml:"title"`
	Link             string `xml:"link"`
	Description      string `xml:"description"`
	HasLocation      bool
	LocationWords    []string
	URLifiedLocation string
	Longitude        float32
	Latitude         float32
}

func findAndFillLocationWords(policeEvents *PoliceEvents) {
	eventsCopy := *policeEvents

	for index, _ := range eventsCopy.Events {
		titleWords, err := filtertitle.FilterTitleWords(eventsCopy.Events[index].Title)
		fmt.Println("Titlewords: ", titleWords)

		if err != nil {
			eventsCopy.Events[index].HasLocation = false
		} else {
			eventsCopy.Events[index].HasLocation = true
			descriptionWords := filterdescription.FilterDescriptionWords(eventsCopy.Events[index].Description)
			fmt.Println("Descwords: ", descriptionWords)
			removeDuplicatesAndCombineLocationWords(titleWords, descriptionWords, &eventsCopy.Events[index].LocationWords)
			AddURLifiedURL(&eventsCopy.Events[index])
		}

	}

	fmt.Println("EventsCopy.Events", eventsCopy.Events[0].LocationWords)

	*policeEvents = eventsCopy
}

func AddURLifiedURL(policeEvent *PoliceEvent) {
	eventCopy := *policeEvent
	str := ""
	for _, word := range eventCopy.LocationWords {
		str += word + " "
	}
	str = url.QueryEscape(str)
	str = strings.TrimSuffix(str, "+")

	eventCopy.URLifiedLocation = str
	*policeEvent = eventCopy
}

func removeDuplicatesAndCombineLocationWords(titleWords []string, descriptionWords []string, locationWords *[]string) {
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

// func moviesearch(wr http.ResponseWriter, re *http.Request) {
// 	//1. lookup omdb-rate (first result only)
// 	//2. lookup rotten-rate (first result only)
// 	//3. create html-page
// 	//4. write html-page with wr

// 	moviename := re.URL.Query().Get("name")

// 	/*if moviename == "" {
// 		wr.Write([]byte("Please enter a valid movie name"))
// 		return
// 	}*/

// 	omovie := omdbquery(moviename)
// 	log.Printf("ImdbMovie: %s, score: %s", omovie.Movietitle, omovie.Imdbscore)

// 	rmovie := rottenquery(moviename)
// 	log.Printf("RottenMovie-score: %d", rmovie.Movies[0].Ratings.Rottenscore)

// 	combinedmoviedata := CombinedMovieData{omovie.Movietitle, omovie.Imdbscore, rmovie.Movies[0].Ratings.Rottenscore}

// 	movietemplate, _ := template.ParseFiles("name.html")
// 	movietemplate.Execute(wr, combinedmoviedata)
// 	//template.Must(template.ParseFiles("name.html")).Execute(wr, combinedmoviedata)

// }

// func omdbquery(moviename string) OmdbMovie {
// 	url := "http://www.omdbapi.com/?t=" + moviename + "&y&plot=short&r=json&tomatoes=true"

// 	omdbresponse, err := http.Get(url)
// 	if err != nil {
// 		log.Println("Error on http.Get: ", err)
// 		return OmdbMovie{}
// 	}

// 	defer omdbresponse.Body.Close()

// 	var omdbmovie OmdbMovie

// 	json.NewDecoder(omdbresponse.Body).Decode(&omdbmovie)

// 	return omdbmovie
// }

// func rottenquery(moviename string) RottenMovie {
// 	url := "http://api.rottentomatoes.com/api/public/v1.0/movies.json?apikey=***REMOVED***&q=" + moviename

// 	rottenresponse, err := http.Get(url)
// 	if err != nil {
// 		log.Println("Error on http.Get: ", err)
// 		return RottenMovie{}
// 	}

// 	defer rottenresponse.Body.Close()

// 	var rottenmovie RottenMovie
// 	json.NewDecoder(rottenresponse.Body).Decode(&rottenmovie)

// 	return rottenmovie
// }

// type OmdbMovie struct {
// 	Movietitle string `json:"Title"`
// 	Imdbscore  string `json:"imdbRating"`
// }

// type RottenMovie struct {
// 	Movies []struct {
// 		MovieTitle string `json:"title"`
// 		Ratings    struct {
// 			Rottenscore int `json:"critics_score"`
// 		} `json:"ratings"`
// 	} `json:"movies"`
// }

// type CombinedMovieData struct {
// 	Movietitle  string
// 	Imdbscore   string
// 	Rottenscore int
// }
