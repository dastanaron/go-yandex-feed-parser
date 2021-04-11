package types

type Offer struct {
	InternalId   string   `xml:"internal-id,attr"`
	Type         string   `xml:"type"`
	PropertyType string   `xml:"property-type"`
	Category     string   `xml:"category"`
	Url          string   `xml:"url"`
	Location     Location `xml:"location"`
}

type Location struct {
	Country         string  `xml:"country"`
	Regon           string  `xml:"region"`
	District        string  `xml:"district"`
	LocalityName    string  `xml:"locality-name"`
	SubLocalityName string  `xml:"sub-locality-name"`
	Address         string  `xml:"address"`
	Direction       string  `xml:"direction"`
	Distance        string  `xml:"distance"`
	Latitude        float64 `xml:"latitude"`
	Longitude       float64 `xml:"longitude"`
	Metro           Metro   `xml:"metro"`
	RailwayStation  string  `xml:"railway-station"`
}

type Metro struct {
	Name            string  `xml:"name"`
	TimeOnFoot      float64 `xml:"time-on-foot"`
	TimeOnTransport float64 `xml:"time-on-transport"`
}
