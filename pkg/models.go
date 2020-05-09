package main

type Destination struct {
	LocationName string `json:"locationName,omitempty"`
}

type TrainService struct {
	STD         string        `json:"std,omitempty"`
	ETD         string        `json:"etd,omitempty"`
	Destination []Destination `json:"destination,omitempty"`
	Platform    string        `json:"platform,omitempty"`
}

type Board struct {
	TrainServices []TrainService `json:"trainServices,omitempty"`
}

type DeparturesByCRSQueryModel struct {
	StationCRS string `json:"stationCRS"`
}
