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
	GetDeparturesByCRS(crs string, boardOptions boardOptions) (*DepartureBoard, error)
	GetArrivalsByCRS(crs string, boardOptions boardOptions) (*ArrivalBoard, error)
}

// Client implements the DepartureBoardIOClient interface by making HTTP requests of the api.departureboard.io JSON API.
type Client struct {
	Client      http.Client
	APIEndpoint url.URL
	APIKey      string
}

func NewClient(apiEndpoint, apiKey string) (client Client, err error) {
	apiURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return client, err
	}
	return Client{
		Client:      http.Client{},
		APIEndpoint: *apiURL,
		APIKey:      apiKey,
	}, nil
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

// FakeClient is an in memory fake that implements the DepartureBoardIOClient interface.
type FakeClient struct {
	Arrivals   map[string][]ArrivalTrainService
	Departures map[string][]DepartureTrainService
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NewFakeClient returns a new FakeClient and an array of valid CRS codes that can be used to make queryies of the client.
func NewFakeClient() (FakeClient, []string) {
	departures := make(map[string][]DepartureTrainService)
	arrivals := make(map[string][]ArrivalTrainService)
	crsCodes := []string{"PAD", "HAY", "NRW", "CBG"}
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
	return FakeClient{
		Departures: departures,
		Arrivals:   arrivals,
	}, crsCodes
}

func (fc *FakeClient) GetDeparturesByCRS(crs string, boardOptions boardOptions) (*DepartureBoard, error) {
	trainServices, ok := fc.Departures[crs]
	if !ok {
		return nil, fmt.Errorf("crs does not exist: %s", crs)
	}
	return &DepartureBoard{
		TrainServices: trainServices,
	}, nil
}

func (fc *FakeClient) GetArrivalsByCRS(crs string, boardOptions boardOptions) (*ArrivalBoard, error) {
	trainServices, ok := fc.Arrivals[crs]
	if !ok {
		return nil, fmt.Errorf("crs does not exist: %s", crs)
	}
	return &ArrivalBoard{
		TrainServices: trainServices,
	}, nil
}
