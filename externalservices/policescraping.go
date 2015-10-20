package externalservices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"poliskarta/structs"
	"sync"
)

func CallPoliceScraping(policeEvent *structs.PoliceEvent, credentials structs.Credentials, wg *sync.WaitGroup) {
	scrapeURL := "https://api.import.io/store/data/3c3e1355-d3c9-4047-bd2e-f86d36af29dc/_query?input/webpage/url="
	apikey := "&_user=" + credentials.Importiouser + "&_apikey=" + credentials.Importiokey

	httpResponse, httperr := http.Get(scrapeURL + policeEvent.PoliceEventURL + apikey)

	if httperr != nil {
		fmt.Println("Importio http-error: " + httperr.Error())
		policeEvent.DescriptionLong = "<N/A>"
	} else {
		defer httpResponse.Body.Close()
		body, ioerr := ioutil.ReadAll(httpResponse.Body)

		if ioerr != nil {
			fmt.Println("Ioutilreadallerror: ", ioerr.Error())
			policeEvent.DescriptionLong = "<N/A>"
		} else {
			var scrapedEvents ScrapedEvents
			unmarshErr := json.Unmarshal(body, &scrapedEvents)

			//For unknown reasons, unmarshal fails some times, might be that the response from
			//police scraping is wrong (200OK instead of a real http-error)
			if unmarshErr != nil {
				fmt.Println("Unmarshal error after police scraping (import.io): " + unmarshErr.Error())
				policeEvent.DescriptionLong = "<N/A>"
			} else {
				//We dont know why, but even though everything seems fine here,
				//policeEvent.DescriptionLong = scrapedEvents.Results[0].Result,
				//sometimes crashes the server. Below is a safety measure.
				for _, result := range scrapedEvents.Results {
					policeEvent.DescriptionLong = result.Result
					break
				}

			}
		}
	}

	defer wg.Done()
}

type ScrapedEvents struct {
	Results []ScrapedEvent `json:"results"`
}
type ScrapedEvent struct {
	Result string `json:"result"`
}
