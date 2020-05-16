package main

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/departureboard-io/departureboard-io-datasource/pkg/departureboardio"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/prometheus/client_golang/prometheus"
)

const metricNamespace = "departureboard_io_datasource"

// DepartureBoardIODataSource handler for departureboard.io API.
type DepartureBoardIODataSource struct {
	DepartureBoardIOClient departureboardio.DepartureBoardIOClient
	APIKey                 string
	APIEndpoint            string
}

// DepartureBoardIOQuery models the query we get from the frontend.
type DepartureBoardIOQuery struct {
	// From the query JSON.
	StationCRS      string `json:"stationCRS"`
	Arrivals        bool   `json:"arrivals"`
	Departures      bool   `json:"departures"`
	FilterCRS       string `json:"filterCRS"`
	ServiceDetails  bool   `json:"serviceDetails"`
	IgnoreTimeRange bool   `json:"ignoreTimeRange"`

	// Not from the query JSON.
	TimeRange backend.TimeRange
}

// NewDataSource creates the departureboard.io datasource
func NewDataSource(mux *http.ServeMux) *DepartureBoardIODataSource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)

	ds := &DepartureBoardIODataSource{
		DepartureBoardIOClient: &departureboardio.Client{
			Client: http.Client{},
		},
	}

	prometheus.MustRegister(queriesTotal)
	return ds
}

// CheckHealth checks if the plugin is running properly
func (ds *DepartureBoardIODataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
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

	if _, err := ds.DepartureBoardIOClient.GetDeparturesByCRS(settings.APIEndpoint, settings.APIKey, "PAD", departureboardio.NewDefaultBoardOptions()); err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Error making request to server"
		return res, nil
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData queries for data.
func (ds *DepartureBoardIODataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()

	settings, err := LoadSettings(req.PluginContext)
	if err != nil {
		return res, errors.New("Unable to load datasource configuration")
	}

	for _, q := range req.Queries {
		dr := ds.handleQuery(settings, q)
		if err != nil {
			return nil, err
		}
		res.Responses[q.RefID] = dr
	}

	return res, nil
}

func (ds *DepartureBoardIODataSource) handleQuery(settings DataSourceSettings, query backend.DataQuery) backend.DataResponse {
	dr := backend.DataResponse{}

	model := DepartureBoardIOQuery{
		TimeRange: query.TimeRange,
	}
	if err := json.Unmarshal(query.JSON, &model); err != nil {
		backend.Logger.Error("Failed to unmarshal query", "query", spew.Sdump(query))
	}

	// Configure boardOptions for departureboard.io queries.
	boardOptions := departureboardio.NewDefaultBoardOptions()
	if model.ServiceDetails {
		boardOptions.ServiceDetails = true
	}
	if model.FilterCRS != "" {
		boardOptions.FilterStation = &model.FilterCRS
	}
	if !model.IgnoreTimeRange {
		timeWindow, timeOffset, err := translateTimeRangeToTimeWindowAndOffset(time.Now(), model.TimeRange)
		if err != nil {
			dr.Error = err
			backend.Logger.Error("Query failed", "model", spew.Sdump(model))
			return dr
		}
		boardOptions.TimeWindow = timeWindow
		boardOptions.TimeOffset = timeOffset
	}

	var frame *data.Frame
	// TODO: is returning multiple frames okay?
	if model.Departures {
		board, err := ds.DepartureBoardIOClient.GetDeparturesByCRS(settings.APIEndpoint, settings.APIKey, model.StationCRS, boardOptions)
		if err != nil {
			dr.Error = err
			backend.Logger.Error("Query failed", "model", spew.Sdump(model))
			return dr
		}
		if model.ServiceDetails {
			frame, err = translateDepartureBoardToFrameWithServiceDetails(model.StationCRS+"Departures", board)
			if err != nil {
				dr.Error = err
				backend.Logger.Error("Failed to translate query response into frame", "response", board)
				return dr
			}
		} else {
			frame, err = translateDepartureBoardToFrame(model.StationCRS+"Departures", board)
			if err != nil {
				dr.Error = err
				backend.Logger.Error("Failed to translate query response into frame", "response", board)
				return dr
			}
		}
		dr.Frames = append(dr.Frames, frame)
	}

	if model.Arrivals {
		board, err := ds.DepartureBoardIOClient.GetArrivalsByCRS(settings.APIEndpoint, settings.APIKey, model.StationCRS, boardOptions)
		if err != nil {
			dr.Error = err
			backend.Logger.Error("Query failed", "model", spew.Sdump(model))
			return dr
		}
		if model.ServiceDetails {
			frame, err = translateArrivalBoardToFrameWithServiceDetails(model.StationCRS+"Arrivals", board)
			if err != nil {
				dr.Error = err
				backend.Logger.Error("Failed to translate query response into frame", "response", board)
				return dr
			}
		} else {
			frame, err = translateArrivalBoardToFrame(model.StationCRS+"Arrivals", board)
			if err != nil {
				dr.Error = err
				backend.Logger.Error("Failed to translate query response into frame", "response", board)
				return dr
			}
		}
		dr.Frames = append(dr.Frames, frame)
	}

	return dr
}

// translateTimeRangeToTimeWindowAndOffset takes the Grafana query time range and translates it to an equivalent representation
// that can be used in a departureboard.io query.
func translateTimeRangeToTimeWindowAndOffset(currentTime time.Time, timeRange backend.TimeRange) (timeWindow, timeOffset int, err error) {
	// Tolerance for clock skew and time between frontend request and backend processing.
	const timeDelta = time.Minute

	if diff := currentTime.Sub(timeRange.To); diff >= 0-timeDelta {
		return timeOffset, timeWindow, errors.New("NationalRail do not provide meaningful historical data, try a time range in the future")
	}
	timeOffset = int(timeRange.From.Sub(currentTime).Minutes())
	timeWindow = int(math.Abs(timeRange.To.Sub(timeRange.From).Minutes()))
	return timeWindow, timeOffset, nil
}

func translateDepartureBoardToFrameWithServiceDetails(name string, board *departureboardio.DepartureBoard) (*data.Frame, error) {
	var destinations, platforms, std, etd, serviceDetails []string
	for _, service := range board.TrainServices {
		std = append(std, service.STD)
		etd = append(etd, service.ETD)
		destinations = append(destinations, service.Destination[0].LocationName)
		platforms = append(platforms, service.Platform)
		callingPoints := []string{}
		// TODO: this doesn't handle a train splitting.
		if len(service.SubsequentCallingPointsList) == 1 {
			for _, cp := range service.SubsequentCallingPointsList[0].SubsequentCallingPoints {
				callingPoints = append(callingPoints, cp.LocationName)
			}
		}
		if len(callingPoints) == 0 {
			serviceDetails = append(serviceDetails, "none")
		} else {
			serviceDetails = append(serviceDetails, strings.Join(callingPoints, ", "))
		}

	}
	return data.NewFrame(name,
		data.NewField("Scheduled", data.Labels{}, std),
		data.NewField("Estimated", data.Labels{}, etd),
		data.NewField("Destination", data.Labels{}, destinations),
		data.NewField("Platform", data.Labels{}, platforms),
		data.NewField("Service Details", data.Labels{}, serviceDetails),
	), nil
}

// translateDepartureBoardToFrame converts a departure board to a data frame.
// TODO: Make a generic board to frame translation.
func translateDepartureBoardToFrame(name string, board *departureboardio.DepartureBoard) (*data.Frame, error) {
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
func translateArrivalBoardToFrame(name string, board *departureboardio.ArrivalBoard) (*data.Frame, error) {
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

func translateArrivalBoardToFrameWithServiceDetails(name string, board *departureboardio.ArrivalBoard) (*data.Frame, error) {
	var origins, platforms, sta, eta, serviceDetails []string
	for _, service := range board.TrainServices {
		sta = append(sta, service.STA)
		eta = append(eta, service.ETA)
		origins = append(origins, service.Origin[0].LocationName)
		platforms = append(platforms, service.Platform)

		callingPoints := []string{}
		if len(service.PreviousCallingPointsList) == 1 {
			for _, cp := range service.PreviousCallingPointsList[0].PreviousCallingPoints {
				callingPoints = append(callingPoints, cp.LocationName)
			}
		}

		if len(callingPoints) == 0 {
			serviceDetails = append(serviceDetails, "none")
		} else {
			serviceDetails = append(serviceDetails, strings.Join(callingPoints, ", "))
		}
	}
	return data.NewFrame(name,
		data.NewField("Scheduled", data.Labels{}, sta),
		data.NewField("Estimated", data.Labels{}, eta),
		data.NewField("Origin", data.Labels{}, origins),
		data.NewField("Platform", data.Labels{}, platforms),
		data.NewField("Service Details", data.Labels{}, serviceDetails),
	), nil
}
