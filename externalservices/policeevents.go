package externalservices

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	ID                    uint32 `json:"ID,string"`
	Title                 string `xml:"title"`
	Link                  string `xml:"link"`
	EventURI              string
	AreaValue             string
	DescriptionShort      string `xml:"description"`
	DescriptionLong       string `json:",omitempty"`
	Time                  string
	EventType             string
	LocationWords         []string `json:",omitempty"`
	CoordinateSearchWords []string `json:",omitempty"`
	Longitude             float32  `json:",omitempty"`
	Latitude              float32  `json:",omitempty"`
}

//Own error type to be able to identify when a faulty event-ID is entered in the URL
type IdNotFoundError struct {
	msg string
}

func (e *IdNotFoundError) Error() string {
	return e.msg
}
