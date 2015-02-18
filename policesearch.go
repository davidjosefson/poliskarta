package main

import "github.com/go-martini/martini"

func main() {
	m := martini.Classic()

	m.Group("/skane", func(r martini.Router) {
		r.Get("/", test)
		r.Get("/wiie", wiie)
		// r.Get("/:id", GetBooks)
	})

	m.Get("/skane/", func() string {
		return "skanesmmmsadfadsfasdf"
	})

	m.Run()

	// http.HandleFunc("/skane/", skane)
	// http.HandleFunc("/skane/0/", skaneSingle)
	// http.HandleFunc("/skane/1/", skaneSingle)
	// http.HandleFunc("/skane/2/", skaneSingle)

	// http.HandleFunc("/stockholm/", stockholm)
	// http.HandleFunc("/stockholm/0/", stockholmSingle)
	// http.HandleFunc("/stockholm/1/", stockholmSingle)

	// //http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css")))) //To find css-files in the css-folder
	// http.ListenAndServe(":9090", nil)
}

func test() string {
	return "groups skane"
}

func wiie() string {
	return "wiee skane"
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
