# poliskarta

Poliskarta is a web API created with [Go](https://golang.org/) to get easy access to coordinates for police events in Sweden.

Event data is fetched and processed from [polisen.se](https://polisen.se/Skane/Aktuellt/Handelser/Handelser-i-hela-landet/?feed=rss) and coordinates are found using [mapquest.se](http://mapquest.com).

Two demo clients were made using AngularJS and Google Maps.

Poliskarta was created as a submission for a university course in Web APIs.

## Running the project

The API needs access to the credentials file **/etc/poliskartaCredentials.json**.

... without it the API will panic.

### API
To start the API run:

    # go run api/main/main.go

.. which will start a web server on port 3000:

 *http://localhost:3000/api/v1/*


### Clients
To start **client-maps-list**:

    # cd client-maps-list/
    # npm install
    # npm run

.. which will start a web server on *http://localhost:7000*

To start **client-big-map**:

    # cd client-big-map/
    # npm install
    # npm run

.. which will start a web server on *http://localhost:8000*
