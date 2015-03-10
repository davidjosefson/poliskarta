package externalservices

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	ID                    uint32 `json:"ID,string"`
	Title                 string `xml:"title"`
	Link                  string `xml:"link"`
	DescriptionShort      string `xml:"description"`
	DescriptionLong       string `json:",omitempty"`
	Time                  string
	EventType             string
	LocationWords         []string `json:",omitempty"`
	CoordinateSearchWords []string `json:",omitempty"`
	Longitude             float32  `json:",omitempty"`
	Latitude              float32  `json:",omitempty"`
}
