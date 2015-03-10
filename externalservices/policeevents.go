package externalservices

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	ID                    uint32 `json: "ID, string"`
	Title                 string `xml:"title"`
	Link                  string `xml:"link"`
	Description           string `xml:"description"`
	Time                  string
	EventType             string
	HasPossibleLocation   bool
	PossibleLocationWords []string
	HasCoordinates        bool
	CoordinateSearchWords []string
	Longitude             float32
	Latitude              float32
}
