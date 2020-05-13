package departureboardio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Client is a client for making requests of the api.departureboard.io JSON API.
type Client struct {
	Client      http.Client
	APIEndpoint url.URL
	APIKey      string
}

// BaseCallingPoint models the common structure in all calling points returned by the departureboard.io API.
type BaseCallingPoint struct {
	// Location name associated calling point.
	LocationName string `json:"locationName,omitempty"`
	// CRS is the Computer Reservation System, three letter name for the calling point.
	CRS string `json:"crs"`
}

// PreviousSubsequentCallingPoint has extra fields that are only present in previous or subsequent calling points
// and not in the origin or destination calling points where that information is part of the train service.
type PreviousSubsequentCallingPoint struct {
	BaseCallingPoint
	// ST is the Scheduled Time that a service is to be at the calling point.
	ST string `json:"st"`
	// ET is the Estimated Time that a service is to be at a calling point.
	// If the ET is equal to the ST, then the ET 'On time', otherwise it is a 24 hour time as a string.
	ET string `json:"et"`
}

// DepartureBoard models the useful elements of a departureboard.io response to a query for departure board.
type DepartureBoard struct {
	TrainServices []DepartureTrainService `json:"trainServices,omitempty"`
}

// ArrivalBoard models the useful elements of a departureboard.io response to a query for arrival board.
type ArrivalBoard struct {
	TrainServices []ArrivalTrainService `json:"trainServices,omitempty"`
}

// BaseTrainService models the common structure of the train services returned by the departureboard.io API.
type BaseTrainService struct {
	// Origin is an array of calling points I think it is only ever one element long and contains the stop from which the train started its journey.
	Origin []BaseCallingPoint `json:"origin,omitempty"`
	// Destination is an array of calling points but I think it is only ever one element long and contains the stop from which the train will end its journey.
	Destination []BaseCallingPoint `json:"destination,omitempty"`
	// Platform is the platform on which the train will stop at for the queried station.
	Platform string `json:"platform,omitempty"`
	// PreviousCallingPointsList is an array of one element 'previousCallingPoints'.
	// It is only present on arrival boards.
	PreviousCallingPointsList []PreviousCallingPointsListElement `json:"previousCallingPointsList,omitempty"`
	// SubsequentCallingPointsList contains arrays of all future calling points.
	// Usually this array is only one element long but can be longer if a service splits.
	// It is only present on departure boards.
	SubsequentCallingPointsList []SubsequentCallingPointsListElement `json:"subsequentCallingPointsList,omitempty"`
}

// DepartureTrainService is a model of a National Rail train service on a departure board.
type DepartureTrainService struct {
	BaseTrainService
	// STD is the Scheduled Time of Departure. It is a 24 hour time as a string.
	STD string `json:"std,omitempty"`
	// ETA is the Estimated Time of Departure. If the ETD is equal to the STD, then the ETD is 'On time', otherwise it is a 24 hour time as a string.
	ETD string `json:"etd,omitempty"`
}

// ArrivalTrainService is a model of a National Rail train service on a arrival board.
type ArrivalTrainService struct {
	BaseTrainService
	// STA is the Scheduled Time of Arrival. It is a 24 hour time as a string.
	STA string `json:"sta,omitempty"`
	// ETA is the Estimated Time of Arrival. If the ETA is equal to the STA, then the ETA is 'On time', otherwise it is a 24 hour time as a string.
	ETA string `json:"eta,omitempty"`
}

type SubsequentCallingPointsListElement struct {
	// SubsequentCallingPoints contains all future calling points.
	SubsequentCallingPoints []PreviousSubsequentCallingPoint `json:"subsequentCallingPoints,omitempty"`
}
type PreviousCallingPointsListElement struct {
	// PreviousCallingPoints contains all scheduled calling points in order from origin to destination including the origin and destination themselves.
	PreviousCallingPoints []PreviousSubsequentCallingPoint `json:"previousCallingPoints,omitempty"`
}

// boardOptions are query parameters that can be set on requests for station arrival or departure boards.
type boardOptions struct {
	// The number of departing services to return.
	// This is a maximum value, less may be returned if there is not a sufficient amount of services running from the selected station within the time window specified.
	NumServices int
	// The time window in minutes to offset the departure information by.
	// For example, a value of 20 will not show services departing within the next 20 minutes.
	TimeOffset int
	// The time window to show services for in minutes.
	// For example, a value of 60 will show services departing the station in the next 60 minutes.
	TimeWindow int
	// Should the response contain information on the calling points for each service?
	// If set to false, calling points will not be returned.
	ServiceDetails bool
	// The CRS (Computer Reservation System) code to filter the results by.
	// For example, performing a lookup for PAD (London Paddington) and setting filterStation to RED (Reading), will only show services departing from Paddington that stop at Reading.
	FilterStation *string
}

// NewBoardOptions returns boardOptions for the provided values. It does not do any validation of its arguments but
// returns an error in case validation is implemented.
func NewBoardOptions(numServices, timeOffset, timeWindow int, serviceDetails bool, filterStation *string) (boardOptions, error) {
	return boardOptions{
		NumServices:    numServices,
		TimeOffset:     timeOffset,
		TimeWindow:     timeWindow,
		ServiceDetails: serviceDetails,
		FilterStation:  filterStation,
	}, nil
}

// NewDefaultBoardOptions returns sensible default options for a Board query.
func NewDefaultBoardOptions() boardOptions {
	return boardOptions{
		NumServices:    20,
		TimeOffset:     20,
		TimeWindow:     60,
		ServiceDetails: false,
		FilterStation:  nil,
	}
}

// getByCRS consolidates the query flow for requests to the getArrivalsByCRS and getDeparturesByCRS endpoints.
func (c *Client) getByCRS(endpoint, crs string, options boardOptions) ([]byte, error) {
	v := url.Values{}
	v.Set("apiKey", c.APIKey)
	v.Set("numServices", strconv.Itoa(options.NumServices))
	v.Set("timeOffset", strconv.Itoa(options.TimeOffset))
	v.Set("timeWindow", strconv.Itoa(options.TimeWindow))
	v.Set("serviceDetails", strconv.FormatBool(options.ServiceDetails))
	if options.FilterStation != nil {
		v.Set("filterStation", *options.FilterStation)
	}

	requestURL := c.APIEndpoint
	requestURL.Path = requestURL.Path + "/" + endpoint + "/" + crs
	requestURL.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Consumer", "DBIO-GRAFANA-PLUGIN")
	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected error from server")
	}

	return body, nil
}

// GetArrivalsByCRS returns a Board of the results from the getDeparturesByCRS endpoint.
func (c *Client) GetArrivalsByCRS(crs string, options boardOptions) (*ArrivalBoard, error) {
	body, err := c.getByCRS("getArrivalsByCRS", crs, options)
	if err != nil {
		return nil, err
	}

	board := &ArrivalBoard{}
	if err := json.Unmarshal(body, board); err != nil {
		return nil, err
	}

	return board, nil
}

// GetDeparturesByCRS returns a Board of the results from the getDeparturesByCRS endpoint.
func (c *Client) GetDeparturesByCRS(crs string, options boardOptions) (*DepartureBoard, error) {
	body, err := c.getByCRS("getDeparturesByCRS", crs, options)
	if err != nil {
		return nil, err
	}

	board := &DepartureBoard{}
	if err := json.Unmarshal(body, board); err != nil {
		return nil, err
	}

	return board, nil
}
