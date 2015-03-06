package externalservices

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	Title                 string `xml:"title"`
	Link                  string `xml:"link"`
	Description           string `xml:"description"`
	HasPossibleLocation   bool
	PossibleLocationWords []string
	HasCoordinates        bool
	CoordinateSearchWords []string
	Longitude             float32
	Latitude              float32
}
