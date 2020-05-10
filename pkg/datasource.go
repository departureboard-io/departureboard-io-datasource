package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/departureboard-io/departureboard-io-datasource/pkg/departureboardio"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/prometheus/client_golang/prometheus"
)

const metricNamespace = "departureboard_io_datasource"

// DepatureBoardIODataSource handler for departureboard.io API.
type DepatureBoardIODataSource struct {
	DepartureBoardIOClient *departureboardio.Client
}

// DepartureBoardIOQuery models the query we get from the frontend.
type DepartureBoardIOQuery struct {
	StationCRS string `json:"stationCRS"`
	Arrivals   bool   `json:"arrivals"`
	Departures bool   `json:"departures"`
}

// NewDataSource creates the departureboard.io datasource
func NewDataSource(mux *http.ServeMux) *DepatureBoardIODataSource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)

	ds := &DepatureBoardIODataSource{
		DepartureBoardIOClient: &departureboardio.Client{
			Client: http.Client{},
		},
	}

	prometheus.MustRegister(queriesTotal)
	return ds
}

// CheckHealth checks if the plugin is running properly
func (ds *DepatureBoardIODataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	// Just checking that the plugin exe is alive and running
	if req.PluginContext.DataSourceInstanceSettings == nil {
		res.Status = backend.HealthStatusOk
		res.Message = "Plugin is Running"
		return res, nil
	}

	settings, err := LoadSettings(req.PluginContext)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid config"
		return res, nil
	}

	apiEndpoint, err := url.Parse(settings.APIEndpoint)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid API endpoint"
	}

	ds.DepartureBoardIOClient.APIEndpoint = *apiEndpoint
	ds.DepartureBoardIOClient.APIKey = settings.APIKey

	ds.DepartureBoardIOClient.GetDeparturesByCRS("PAD", departureboardio.NewDefaultBoardOptions())
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Error making request to server"
		return res, nil
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData queries for data.
func (ds *DepatureBoardIODataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()

	if ds.DepartureBoardIOClient.APIKey == "" {
		settings, err := LoadSettings(req.PluginContext)
		if err != nil {
			return res, errors.New("Unable to load datasource configuration")
		}

		apiEndpoint, err := url.Parse(settings.APIEndpoint)
		if err != nil {
			return res, errors.New("Unable to load datasource configuration")
		}

		ds.DepartureBoardIOClient.APIEndpoint = *apiEndpoint
		ds.DepartureBoardIOClient.APIKey = settings.APIKey
	}

	for _, q := range req.Queries {
		dr := backend.DataResponse{}
		model := &DepartureBoardIOQuery{}
		if err := json.Unmarshal(q.JSON, model); err != nil {
			backend.Logger.Error("Query failed", "model", spew.Sdump(model))
		}

		boardOptions := departureboardio.NewDefaultBoardOptions()
		// TODO: is returning multiple frames okay?
		if model.Departures {
			board, err := ds.DepartureBoardIOClient.GetDeparturesByCRS(model.StationCRS, boardOptions)
			if err != nil {
				dr.Error = err
				return res, nil
			}
			frame, err := translateDepartureBoardToFrame(model.StationCRS+"Departures", board)
			if err != nil {
				dr.Error = err
				return res, nil
			}
			dr.Frames = append(dr.Frames, frame)
		}

		if model.Arrivals {
			board, err := ds.DepartureBoardIOClient.GetArrivalsByCRS(model.StationCRS, boardOptions)
			if err != nil {
				dr.Error = err
				return res, nil
			}
			frame, err := translateArrivalBoardToFrame(model.StationCRS+"Arrivals", board)
			if err != nil {
				dr.Error = err
				return res, nil
			}
			dr.Frames = append(dr.Frames, frame)
		}

		res.Responses[q.RefID] = dr
	}

	return res, nil
}

// translateDepartureBoardToFrame converts a departure board to a data frame.
// TODO: Make a generic board to frame translation.
func translateDepartureBoardToFrame(name string, board *departureboardio.Board) (*data.Frame, error) {
	var destinations, platforms, std, etd []string
	for _, service := range board.TrainServices {
		std = append(std, service.STD)
		etd = append(etd, service.ETD)
		destinations = append(destinations, service.Destination[0].LocationName)
		platforms = append(platforms, service.Platform)
	}
	return data.NewFrame(name,
		data.NewField("Scheduled", data.Labels{}, std),
		data.NewField("Estimated", data.Labels{}, etd),
		data.NewField("Destination", data.Labels{}, destinations),
		data.NewField("Platform", data.Labels{}, platforms),
	), nil
}

// translateArrivalBoardToFrame converts a departure board to a data frame.
// TODO: Make a generic board to frame translation.
func translateArrivalBoardToFrame(name string, board *departureboardio.Board) (*data.Frame, error) {
	var origins, platforms, sta, eta []string
	for _, service := range board.TrainServices {
		sta = append(sta, service.STA)
		eta = append(eta, service.ETA)
		origins = append(origins, service.Origin[0].LocationName)
		platforms = append(platforms, service.Platform)
	}
	return data.NewFrame(name,
		data.NewField("Scheduled", data.Labels{}, sta),
		data.NewField("Estimated", data.Labels{}, eta),
		data.NewField("Origin", data.Labels{}, origins),
		data.NewField("Platform", data.Labels{}, platforms),
	), nil
}
