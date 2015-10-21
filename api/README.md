# poliskarta: API
A web API written in Go for easy access to swedish police events and coordinates for where they took place.

## Prerequisites
* Download and install [Go](https://golang.org/dl/)
* The API needs access to the credentials file **/etc/poliskartaCredentials.json**.

The credentials file consists of API-keys to MapQuest and Import.io.

## Start the API

    # go run main/main.go

.. which will start a web server on port 3000.

You can then access the api here:

    http://localhost:3000/api/v1/


## Endpoints

* All responses is in JSON
* The endpoints only accept GET-requests

### /areas
Provides array of all the areas the swedish police publishes events for.

#### Example response
```json
{  
   "areas":[  
      {  
         "name":"Blekinge",
         "value":"blekinge",
         "latitude":56.283333,
         "longitude":15.116667,
         "zoomlevel":8,
         "links":[  
            {  
               "rel":"self",
               "href":"http://localhost:3000/api/v1/areas/blekinge"
            },
            {  
               "rel":"origin",
               "href":"https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss"
            }
         ]
      },
      ...
   ]
}
```

### /areas/<area>
Provides information of the requested area together with the latest events that occured in the area published by the swedish police.

Observe that the events in the response for this endpoint DOES NOT include coordinates, identified words for places or the extended information about the event. To get this information one has to make a GET request for /areas/<area>/<event-ID>.

#### Optional parameter: ?limit=35
Example usage:

    GET http://localhost:3000/api/v1/areas/blekinge?limit=15


* Used to limit the number of events the response should include.
* Valid values: **1-50** (values above 50 will result in 50 events)

#### Example response
A GET request for http://localhost:3000:/api/v1/areas/blekinge will generate the following:

```json
{  
   "name":"Blekinge",
   "value":"blekinge",
   "latitude":56.283333,
   "longitude":15.116667,
   "zoomlevel":8,
   "events":[  
      {  
         "id":"1416040185",
         "title":"2015-03-19 06:57, Inbrott, Karlshamn",
         "time":"2015-03-19 06:57",
         "eventType":"Inbrott",
         "descriptionShort":"Inbrott på lager, Kungsgatan.",
         "location":{  
            "words":[  
               "Kungsgatan",
               "Karlshamn"
            ]
         },
         "links":[  
            {  
               "rel":"self",
               "href":"http://localhost:3000/api/v1/areas/blekinge/1416040185"
            },
            {  
               "rel":"origin",
               "href":"http://polisen.se/Halland/Aktuellt/Handelser/Blekinge/2015-03-19-0657-Inbrott-Karlshamn/"
            }
         ]
      },
      ...
   ],
   "links":[  
      {  
         "rel":"self",
         "href":"http://localhost:3000/api/v1/areas/blekinge"
      },
      {  
         "rel":"origin",
         "href":"https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss"
      }
   ]
}
```

### /areas/<area>/<event-id>
Provides detailed information about a certain event.

A GET request for http://localhost:3000/api/v1/areas/blekinge/1416040185 will generate the following:

```json
{  
   "id":"1416040185",
   "title":"2015-03-19 06:57, Inbrott, Karlshamn",
   "time":"2015-03-19 06:57",
   "eventType":"Inbrott",
   "descriptionShort":"Inbrott på lager, Kungsgatan.",
   "descriptionLong":"Inbrott konstateras i flera lager på Kungsgatan där lagren uppges tillhöra flera olika butiker. Vad som tillgripits är inledningsvis okänt.",
   "area":{  
      "name":"Blekinge",
      "value":"blekinge",
      "links":[  
         {  
            "rel":"self",
            "href":"http://localhost:3000/api/v1/areas/blekinge"
         },
         {  
            "rel":"origin",
            "href":"https://polisen.se/Halland/Aktuellt/RSS/Lokal-RSS---Handelser/Lokala-RSS-listor1/Handelser-RSS---Blekinge/?feed=rss"
         }
      ]
   },
   "location":{  
      "words":[  
         "Kungsgatan",
         "Karlshamn"
      ],
      "searchWords":[  
         "Kungsgatan",
         "Karlshamn"
      ],
      "longitude":14.864794,
      "latitude":56.169296
   },
   "links":[  
      {  
         "rel":"self",
         "href":"http://localhost:3000/api/v1/areas/blekinge/1416040185"
      },
      {  
         "rel":"origin",
         "href":"http://polisen.se/Halland/Aktuellt/Handelser/Blekinge/2015-03-19-0657-Inbrott-Karlshamn/"
      }
   ]
}
```
