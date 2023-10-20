package models

type Route struct {
	AirportFrom       string `json:"airportFrom"`
	AirportTo         string `json:"airportTo"`
	ConnectingAirport string `json:"connectingAirport"`
	NewRoute          bool   `json:"newRoute"`
	SeasonalRoute     bool   `json:"seasonalRoute"`
	Operator          string `json:"operator"`
	Group             string `json:"group"`
}
