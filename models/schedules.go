package models

type Schedule struct {
	Month int   `json:"month"`
	Days  []Day `json:"days"`
}

type Day struct {
	Day     int      `json:"day"`
	Flights []Flight `json:"flights"`
}

type Flight struct {
	Number        string `json:"number,omitempty"`
	DepartureTime string `json:"departureTime,omitempty"`
	ArrivalTime   string `json:"arrivalTime,omitempty"`
	Stops         int    `json:"stops"`
	Legs          []Leg  `json:"leg,omitempty"`
}

type Leg struct {
	DepartureAirport  string `json:"departureAirport"`
	ArrivalAirport    string `json:"arrivalAirport"`
	DepartureDateTime string `json:"departureDateTime"`
	ArrivalDateTime   string `json:"arrivalDateTime"`
}

type InterconnectedFlight struct {
	Stops int   `json:"stops"`
	Legs  []Leg `json:"legs"`
}

type FlightDetail struct {
	Day              int    `json:"day"`
	Month            int    `json:"month"`
	CarrierCode      string `json:"carrierCode"`
	Number           string `json:"number"`
	DepartureTime    string `json:"departureTime"`
	ArrivalTime      string `json:"arrivalTime"`
	DepartureAirport string `json:"departureAirport"`
	ArrivalAirport   string `json:"arrivalAirport"`
}
