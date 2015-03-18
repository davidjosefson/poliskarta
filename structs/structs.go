package structs

type AreasStruct struct {
	Areas []Area `json:"areas"`
}
type Area struct {
	Name            string  `json:"name"`
	Value           string  `json:"value"`
	RssURL          string  `json:"-"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	GoogleZoomLevel int     `json:"zoomlevel"`
	Links           []Link  `json:"links"`
}

type PoliceEvents struct {
	Events []PoliceEvent `xml:"channel>item"`
}

type PoliceEvent struct {
	ID               uint32          `json:"id,string"`
	Title            string          `xml:"title" json:"title"`
	Time             string          `json:"time"`
	EventType        string          `json:"eventType"`
	DescriptionShort string          `xml:"description" json:"descriptionShort"`
	DescriptionLong  string          `json:"descriptionLong ,omitempty"`
	Area             PoliceEventArea `json:"area"`
	Location         LocationInfo    `json:"location"`
	PoliceEventURL   string          `xml:"link" json:"-"`
	Links            []Link          `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type LocationInfo struct {
	Words       []string `json:"words,omitempty"`
	SearchWords []string `json:"searchWords,omitempty"`
	Longitude   float32  `json:"longitude,omitempty"`
	Latitude    float32  `json:"latitude,omitempty"`
}

type PoliceEventArea struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Links []Link `json:"links"`
}

//Own error type to be able to identify when a faulty event-ID is entered in the URL
type IdNotFoundError struct {
	Msg string
}

func (e *IdNotFoundError) Error() string {
	return e.Msg
}

type GeoLocation struct {
	Locations []Location `xml:"results>result"`
	// ThumbMaps string     `xml:"options>thumbMaps"`
}

type Location struct {
	LocationAlternatives []LocationAlternative `xml:"locations>location"`
}

type LocationAlternative struct {
	Quality   string  `xml:"geocodeQuality"`
	Latitude  float32 `xml:"displayLatLng>latLng>lat"`
	Longitude float32 `xml:"displayLatLng>latLng>lng"`
}
