package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"
)

// var placesSlice = []string{
// 	"blekinge", "dalarna", "gotland", "gavleborg", "halland", "jamtland",
// 	"jonkoping", "kalmar", "kronoberg", "norrbotten", "skane", "stockholm",
// 	"sodermanland", "uppsala", "varmland", "vasterbotten", "vasternorrland",
// 	"vastmanland", "vastragotaland", "orebro", "ostergotland"}

var places = map[string]string{
	"blekinge":       "http://www.polisen.se/rss-handelser-blekinge",
	"dalarna":        "http://www.polisen.se/rss-handelser-dalarna",
	"gotland":        "http://www.polisen.se/rss-handelser-gotland",
	"gavleborg":      "http://www.polisen.se/rss-handelser-gavleborg",
	"halland":        "http://www.polisen.se/rss-handelser-halland",
	"jamtland":       "http://www.polisen.se/rss-handelser-jamtland",
	"jonkoping":      "http://www.polisen.se/rss-handelser-jonkoping",
	"kalmar":         "http://www.polisen.se/rss-handelser-kalmar",
	"kronoberg":      "http://www.polisen.se/rss-handelser-kronoberg",
	"norrbotten":     "http://www.polisen.se/rss-handelser-norrbotten",
	"skane":          "http://www.polisen.se/rss-handelser-skane",
	"stockholm":      "http://www.polisen.se/rss-handelser-stockholm",
	"sodermanland":   "http://www.polisen.se/rss-handelser-sodermanland",
	"uppsala":        "http://www.polisen.se/rss-handelser-uppsala",
	"varmland":       "http://www.polisen.se/rss-handelser-varmland",
	"vasterbotten":   "http://www.polisen.se/rss-handelser-vasterbotten",
	"vasternorrland": "http://www.polisen.se/rss-handelser-vasternorrland",
	"vastmanland":    "http://www.polisen.se/rss-handelser-vastmanland",
	"vastragotaland": "http://www.polisen.se/rss-handelser-vastragotaland",
	"orebro":         "http://www.polisen.se/rss-handelser-orebro",
	"ostergotland":   "http://www.polisen.se/rss-handelser-ostergotland",
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
		res.Header().Add("Content-type", "application/json")
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

//SKAPAR INTE NÅGON JSON, HÄMTAR BARA XML FRÅN POLISEN OCH RETURNERAR!
func callExternalServicesAndCreateJson(place string) string {
	response, _ := http.Get(places[place])
	str, _ := ioutil.ReadAll(response.Body)
	return string(str)
}

// func skane(wr http.ResponseWriter, re *http.Request) {
// 	/*
// 		1. Fixa så att man kan anropa EN metod för alla län, en array med alla unika urls till polisens
// 		2. Hämta polis-RSS och mappa till struct för Länet
// 		3. Fixa en metod som kan ta reda på platsen namn (stad, by osv) för att söka i Google Maps
// 		4. Gör ett anrop till Google Maps för varje platsnamn och få tillbaka koordinater/platsnamn
// 		5. Returnera en lång lista med händelser och platskoordinater

// 		6. Ifall man går in på skane/1/ ska enbart PoliceEvent[0] för det länet returneras
// 	*/

// 	policeresponse, _ := http.Get("https://polisen.se/Gotlands_lan/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Gotland/?feed=rss")

// 	defer policeresponse.Body.Close()

// 	var channel Channel

// 	xml.NewDecoder(policeresponse.Body).Decode(&channel)

// 	fmt.Println(channel.Items)
// }

// type Foobar struct {
// 	PoliceEvents []PoliceEvent `xml:"channel>item"`
// }

// type PoliceEvent struct {
// 	Title       string `xml:"title"`
// 	Link        string `xml:"link"`
// 	Description string `xml:"description"`
// 	PubDate     string `xml:"pubDate"`
// }

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
