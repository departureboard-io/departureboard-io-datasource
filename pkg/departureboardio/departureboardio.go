package departureboardio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type DepartureBoardIOClient interface {
	GetDeparturesByCRS(apiEndpoint, apiKey, crs string, boardOptions boardOptions) (*DepartureBoard, error)
	GetArrivalsByCRS(apiEndpoint, apiKey, crs string, boardOptions boardOptions) (*ArrivalBoard, error)
}

// Client implements the DepartureBoardIOClient interface by making HTTP requests of the api.departureboard.io JSON API.
type Client struct {
	Client http.Client
}

func NewClient(apiEndpoint, apiKey string) (client Client, err error) {
	if err != nil {
		return client, err
	}
	return Client{Client: http.Client{}}, nil
}

// getByCRS consolidates the query flow for requests to the getArrivalsByCRS and getDeparturesByCRS endpoints.
func (c *Client) getByCRS(apiEndpoint, apiKey, queryEndpoint, crs string, options boardOptions) ([]byte, error) {
	v := url.Values{}
	v.Set("apiKey", apiKey)
	v.Set("numServices", strconv.Itoa(options.NumServices))
	v.Set("timeOffset", strconv.Itoa(options.TimeOffset))
	v.Set("timeWindow", strconv.Itoa(options.TimeWindow))
	v.Set("serviceDetails", strconv.FormatBool(options.ServiceDetails))
	if options.FilterStation != nil {
		v.Set("filterStation", *options.FilterStation)
	}

	requestURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse api endpoint as URL: %v", err)
	}
	requestURL.Path = requestURL.Path + "/" + queryEndpoint + "/" + crs
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
func (c Client) GetArrivalsByCRS(apiEndpoint, apiKey, crs string, options boardOptions) (*ArrivalBoard, error) {
	body, err := c.getByCRS(apiEndpoint, apiKey, "getArrivalsByCRS", crs, options)
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
func (c Client) GetDeparturesByCRS(apiEndpoint, apiKey, crs string, options boardOptions) (*DepartureBoard, error) {
	body, err := c.getByCRS(apiEndpoint, apiKey, "getDeparturesByCRS", crs, options)
	if err != nil {
		return nil, err
	}

	board := &DepartureBoard{}
	if err := json.Unmarshal(body, board); err != nil {
		return nil, err
	}

	return board, nil
}

// FakeClient is an in memory fake that implements the DepartureBoardIOClient interface.
type FakeClient struct {
	APIEndpoint url.URL
	APIKey      string
	Arrivals    map[string][]ArrivalTrainService
	Departures  map[string][]DepartureTrainService
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NewFakeClient returns a new FakeClient with arrivals and departure train services generated for all provided crsCodes.
func NewFakeClient(crsCodes []string) *FakeClient {
	departures := make(map[string][]DepartureTrainService)
	arrivals := make(map[string][]ArrivalTrainService)
	for _, crs := range crsCodes {
		departures[crs] = []DepartureTrainService{
			{
				BaseTrainService: BaseTrainService{
					Platform:    "12A",
					Destination: []BaseCallingPoint{{LocationName: crs}},
					Origin:      []BaseCallingPoint{{LocationName: reverseString(crs)}},
				},
				STD: "13:10",
				ETD: "On time",
				SubsequentCallingPointsList: []SubsequentCallingPointsListElement{
					SubsequentCallingPointsListElement{
						[]PreviousSubsequentCallingPoint{
							{
								BaseCallingPoint: BaseCallingPoint{
									LocationName: crs,
								},
								ST: "14:00",
								ET: "14:01",
							},
						},
					},
				},
			},
		}
		arrivals[crs] = []ArrivalTrainService{
			{
				BaseTrainService: BaseTrainService{
					Platform:    "12A",
					Origin:      []BaseCallingPoint{{LocationName: crs}},
					Destination: []BaseCallingPoint{{LocationName: reverseString(crs)}},
				},
				STA: "13:10",
				ETA: "On time",
				PreviousCallingPointsList: []PreviousCallingPointsListElement{
					PreviousCallingPointsListElement{
						[]PreviousSubsequentCallingPoint{
							{
								BaseCallingPoint: BaseCallingPoint{
									LocationName: crs,
								},
								ST: "14:00",
								ET: "14:01",
							},
						},
					},
				},
			},
		}
	}
	return &FakeClient{
		Departures: departures,
		Arrivals:   arrivals,
	}
}

func (fc FakeClient) GetDeparturesByCRS(apiKey, apiEndpoint, crs string, boardOptions boardOptions) (*DepartureBoard, error) {
	trainServices, ok := fc.Departures[crs]
	if !ok {
		return nil, fmt.Errorf("crs does not exist: %s", crs)
	}
	return &DepartureBoard{
		TrainServices: trainServices,
	}, nil
}

func (fc FakeClient) GetArrivalsByCRS(apiKey, apiEndpoint, crs string, boardOptions boardOptions) (*ArrivalBoard, error) {
	trainServices, ok := fc.Arrivals[crs]
	if !ok {
		return nil, fmt.Errorf("crs does not exist: %s", crs)
	}
	return &ArrivalBoard{
		TrainServices: trainServices,
	}, nil
}
