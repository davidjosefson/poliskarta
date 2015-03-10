package externalservices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func CallPoliceScraping(policeEvent *PoliceEvent, wg *sync.WaitGroup) {
	copyEvent := *policeEvent
	scrapeURL := "https://api.import.io/store/data/3c3e1355-d3c9-4047-bd2e-f86d36af29dc/_query?input/webpage/url="
	apikey := "&_user=***REMOVED***&_apikey=***REMOVED***"

	httpResult, httperr := http.Get(scrapeURL + copyEvent.Link + apikey)

	if httperr != nil {
		fmt.Println("Importio http-error: " + httperr.Error())
	} else {
		body, ioerr := ioutil.ReadAll(httpResult.Body)
		if ioerr != nil {
			fmt.Println("Ioutilreadallerror: ", ioerr.Error())
		} else {
			fmt.Println(string(body))

			var scrapedEvents ScrapedEvents
			json.Unmarshal(body, &scrapedEvents)

			copyEvent.DescriptionLong = scrapedEvents.Results[0].Result

		}
	}
	*policeEvent = copyEvent
	defer wg.Done()
}

type ScrapedEvents struct {
	Results []ScrapedEvent `json:"results"`
}
type ScrapedEvent struct {
	Result string `json:"result"`
}
