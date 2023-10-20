package models

type RequestData struct {
	Departure         string `json:"departure"`
	DepartureDateTime string `json:"departureDateTime"`
	Arrival           string `json:"arrival"`
	ArrivalDateTime   string `json:"arrivalDateTime"`
}
